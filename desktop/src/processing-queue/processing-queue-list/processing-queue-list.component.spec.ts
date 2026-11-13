import { NO_ERRORS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute, convertToParamMap } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { of, throwError } from "rxjs";
import { Group, Permission, ProcessingQueueItem, ProcessingQueueService } from "../../open-api";
import { PipesModule } from "../../pipes";
import { SnackbarService } from "../../services/snackbar.service";
import { AuthState, GroupState, SetGroups, SetPermissions, SetSelectedGroupId } from "../../store";
import { SystemSettingsState } from "../../store/system-settings.state";
import { ProcessingQueueCountService } from "../processing-queue-count.service";
import { PROCESSING_POLL_INTERVAL_MS } from "../processing-queue.utils";
import { ProcessingQueueListComponent } from "./processing-queue-list.component";

describe("ProcessingQueueListComponent", () => {
  let fixture: ComponentFixture<ProcessingQueueListComponent>;
  let component: ProcessingQueueListComponent;

  const api = {
    getProcessingQueueForGroup: jest.fn(),
    getProcessingQueueTaskSourceFile: jest.fn(),
    downloadProcessingQueueTaskSourceFile: jest.fn(),
    retryProcessingQueueTask: jest.fn(),
  };
  const countService = { report: jest.fn(), refresh: jest.fn() };
  const dialog = { open: jest.fn() };
  const snackbar = { success: jest.fn(), error: jest.fn() };

  const home = { id: 5, name: "Home", isAllGroup: false } as Group;
  const work = { id: 6, name: "Work", isAllGroup: false } as Group;
  const all = { id: 9, name: "All", isAllGroup: true } as Group;

  const item = (overrides: Partial<ProcessingQueueItem> = {}): ProcessingQueueItem => ({
    asynqTaskId: "t1",
    queue: "quick_scan",
    state: "QUEUED",
    taskType: "QUICK_SCAN",
    groupId: 5,
    fileName: "receipt.png",
    retryCount: 0,
    maxRetries: 3,
    lastError: "",
    hasSourceFile: true,
    canBeRestarted: false,
    canRetryNow: false,
    ...overrides,
  });

  const respond = (items: ProcessingQueueItem[], truncated = false) =>
    api.getProcessingQueueForGroup.mockReturnValue(of({ items, truncated }));

  const setUp = async (groupId = 5) => {
    await TestBed.configureTestingModule({
      declarations: [ProcessingQueueListComponent],
      imports: [PipesModule, NgxsModule.forRoot([AuthState, GroupState, SystemSettingsState])],
      providers: [
        provideZonelessChangeDetection(),
        { provide: ProcessingQueueService, useValue: api },
        { provide: ProcessingQueueCountService, useValue: countService },
        { provide: MatDialog, useValue: dialog },
        { provide: SnackbarService, useValue: snackbar },
        {
          provide: ActivatedRoute,
          useValue: {
            paramMap: of(convertToParamMap({ groupId: String(groupId) })),
            snapshot: { paramMap: convertToParamMap({ groupId: String(groupId) }) },
          },
        },
      ],
      schemas: [NO_ERRORS_SCHEMA],
    }).compileComponents();

    const store = TestBed.inject(Store);
    store.dispatch(new SetGroups([all, home, work]));
    store.dispatch(new SetSelectedGroupId(String(groupId)));
    store.dispatch(new SetPermissions([], { 5: [Permission.GroupActivitiesRead] }));

    fixture = TestBed.createComponent(ProcessingQueueListComponent);
    component = fixture.componentInstance;
    fixture.autoDetectChanges();
    await fixture.whenStable();
  };

  const byTestId = (id: string): HTMLElement | null =>
    fixture.nativeElement.querySelector(`[data-testid='${id}']`);

  const setHidden = (hidden: boolean) => {
    Object.defineProperty(document, "hidden", { configurable: true, get: () => hidden });
  };

  beforeEach(() => {
    jest.resetAllMocks();
    setHidden(false);
    localStorage.clear();
  });

  afterEach(() => {
    jest.useRealTimers();
    setHidden(false);
  });

  it("loads the route group and reports the count to the badge", async () => {
    respond([item(), item({ asynqTaskId: "t2", state: "FAILED" })]);
    await setUp();

    expect(api.getProcessingQueueForGroup).toHaveBeenCalledWith(5);
    expect(component.items().length).toBe(2);
    expect(countService.report).toHaveBeenCalledWith(5, 2);
  });

  it("shows the empty message and no table when nothing is queued", async () => {
    respond([]);
    await setUp();

    expect(byTestId("processing-queue-empty")?.textContent).toContain("Nothing is processing.");
    expect(fixture.nativeElement.querySelector("app-table")).toBeNull();
  });

  it("notes a truncated list", async () => {
    respond([item()], true);
    await setUp();

    expect(byTestId("processing-queue-truncated")).not.toBeNull();
  });

  it("keeps an unknown state in flight and renders its raw code", async () => {
    respond([item({ state: "PAUSED", taskType: "FAX_UPLOAD" })]);
    await setUp();

    expect(component.polling()).toBe(true);
    expect(component.stateLabel("PAUSED")).toBe("PAUSED");
    expect(component.typeLabel("FAX_UPLOAD")).toBe("FAX_UPLOAD");
  });

  it("adds a Group column on the All group only", async () => {
    respond([item()]);
    await setUp(9);
    expect(component.displayedColumns()).toContain("group");
  });

  it("has no Group column on a real group", async () => {
    respond([item()]);
    await setUp(5);
    expect(component.displayedColumns()).not.toContain("group");
  });

  describe("row actions", () => {
    it("offers preview and download only with a source file", async () => {
      respond([item({ hasSourceFile: false })]);
      await setUp();
      const template = component["actionsCell"]();
      const view = template.createEmbeddedView({ element: component.items()[0] });
      view.detectChanges();
      expect(view.rootNodes.some((n: HTMLElement) => n.querySelector?.("[data-testid='processing-queue-preview']"))).toBe(false);
    });

    it("words the retry button by state", async () => {
      respond([]);
      await setUp();
      expect(component.retryText(item({ state: "RETRYING", canRetryNow: true }))).toBe("Retry now");
      expect(component.retryText(item({ state: "FAILED", canBeRestarted: true }))).toBe("Resubmit");
      expect(component.retryText(item({ state: "QUEUED" }))).toBeUndefined();
    });

    it("retries a RETRYING item at once, without a confirmation, then reloads", async () => {
      const retrying = item({ state: "RETRYING", canRetryNow: true });
      respond([retrying]);
      api.retryProcessingQueueTask.mockReturnValue(of({}));
      await setUp();
      api.getProcessingQueueForGroup.mockClear();

      component.onRetry(retrying);

      expect(dialog.open).not.toHaveBeenCalled();
      expect(api.retryProcessingQueueTask).toHaveBeenCalledWith("quick_scan", "t1");
      expect(api.getProcessingQueueForGroup).toHaveBeenCalledTimes(1);
    });

    it("confirms before resubmitting a FAILED item, and does nothing if declined", async () => {
      const failed = item({ state: "FAILED", canBeRestarted: true });
      respond([failed]);
      await setUp();
      const afterClosed = jest.fn().mockReturnValue(of(false));
      dialog.open.mockReturnValue({ componentInstance: {}, afterClosed });

      component.onRetry(failed);

      expect(dialog.open).toHaveBeenCalled();
      expect(api.retryProcessingQueueTask).not.toHaveBeenCalled();

      afterClosed.mockReturnValue(of(true));
      api.retryProcessingQueueTask.mockReturnValue(of({}));
      component.onRetry(failed);
      expect(api.retryProcessingQueueTask).toHaveBeenCalledWith("quick_scan", "t1");
    });

    it("reloads when a retry is refused, since the task may have moved on", async () => {
      const retrying = item({ state: "RETRYING", canRetryNow: true });
      respond([retrying]);
      api.retryProcessingQueueTask.mockReturnValue(throwError(() => new Error("400")));
      await setUp();
      api.getProcessingQueueForGroup.mockClear();

      component.onRetry(retrying);

      expect(api.getProcessingQueueForGroup).toHaveBeenCalledTimes(1);
      expect(component.isRetrying(retrying)).toBe(false);
    });

    it("previews through the queue route, not the activity route", async () => {
      respond([item()]);
      api.getProcessingQueueTaskSourceFile.mockReturnValue(of({ name: "receipt.png", encodedImage: "data:x" }));
      await setUp();

      component.onPreview(component.items()[0]);

      expect(api.getProcessingQueueTaskSourceFile).toHaveBeenCalledWith("quick_scan", "t1");
      expect(dialog.open).toHaveBeenCalled();
    });
  });

  describe("polling", () => {
    beforeEach(() => jest.useFakeTimers());

    it("refreshes every interval while something is in flight", async () => {
      respond([item({ state: "PROCESSING" })]);
      await setUp();
      expect(api.getProcessingQueueForGroup).toHaveBeenCalledTimes(1);

      jest.advanceTimersByTime(PROCESSING_POLL_INTERVAL_MS);
      expect(api.getProcessingQueueForGroup).toHaveBeenCalledTimes(2);
      jest.advanceTimersByTime(PROCESSING_POLL_INTERVAL_MS);
      expect(api.getProcessingQueueForGroup).toHaveBeenCalledTimes(3);
    });

    it("stops once nothing is in flight (failed items do not poll)", async () => {
      respond([item({ state: "FAILED" })]);
      await setUp();

      jest.advanceTimersByTime(PROCESSING_POLL_INTERVAL_MS * 3);
      expect(api.getProcessingQueueForGroup).toHaveBeenCalledTimes(1);
    });

    it("does not poll while the tab is hidden, and reloads when it returns", async () => {
      setHidden(true);
      respond([item({ state: "QUEUED" })]);
      await setUp();

      jest.advanceTimersByTime(PROCESSING_POLL_INTERVAL_MS * 3);
      expect(api.getProcessingQueueForGroup).toHaveBeenCalledTimes(1);

      setHidden(false);
      document.dispatchEvent(new Event("visibilitychange"));
      expect(api.getProcessingQueueForGroup).toHaveBeenCalledTimes(2);
    });

    it("stops polling when the page is left", async () => {
      respond([item({ state: "QUEUED" })]);
      await setUp();

      fixture.destroy();
      jest.advanceTimersByTime(PROCESSING_POLL_INTERVAL_MS * 3);
      expect(api.getProcessingQueueForGroup).toHaveBeenCalledTimes(1);
    });

    it("stops polling after a failed load and says so", async () => {
      api.getProcessingQueueForGroup.mockReturnValue(throwError(() => new Error("down")));
      await setUp();

      jest.advanceTimersByTime(PROCESSING_POLL_INTERVAL_MS * 3);
      expect(api.getProcessingQueueForGroup).toHaveBeenCalledTimes(1);
      expect(byTestId("processing-queue-error")).not.toBeNull();
    });
  });
});
