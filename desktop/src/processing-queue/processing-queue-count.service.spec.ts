import { provideZonelessChangeDetection } from "@angular/core";
import { TestBed } from "@angular/core/testing";
import { NavigationEnd, Router } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { of, Subject, throwError } from "rxjs";
import { Group, Permission, ProcessingQueueService } from "../open-api";
import { AuthState, GroupState, SetAuthState, SetGroups, SetPermissions, SetSelectedGroupId } from "../store";
import { ProcessingQueueCountService } from "./processing-queue-count.service";

describe("ProcessingQueueCountService", () => {
  const api = { getProcessingQueueForGroup: jest.fn() };
  let events: Subject<unknown>;
  const real = { id: 5, name: "Home", isAllGroup: false } as Group;
  const all = { id: 9, name: "All", isAllGroup: true } as Group;
  const list = (n: number) => ({ items: new Array(n).fill({}), truncated: false });

  let service: ProcessingQueueCountService;
  let store: Store;

  const setUp = (permissions: { [id: number]: string[] }) => {
    store.dispatch(new SetAuthState({ exp: Math.floor(Date.now() / 1000) + 3600 } as any));
    store.dispatch(new SetGroups([all, real]));
    store.dispatch(new SetPermissions([], permissions));
    store.dispatch(new SetSelectedGroupId("5"));
    TestBed.tick();
  };

  beforeEach(() => {
    jest.resetAllMocks();
    localStorage.clear();
    events = new Subject<unknown>();
    api.getProcessingQueueForGroup.mockReturnValue(of(list(3)));
    TestBed.configureTestingModule({
      imports: [NgxsModule.forRoot([AuthState, GroupState])],
      providers: [
        provideZonelessChangeDetection(),
        { provide: ProcessingQueueService, useValue: api },
        { provide: Router, useValue: { events } },
      ],
    });
    store = TestBed.inject(Store);
    service = TestBed.inject(ProcessingQueueCountService);
  });

  it("fetches the selected group's count once the viewer may read the queue", () => {
    setUp({ 5: [Permission.GroupActivitiesRead] });

    expect(api.getProcessingQueueForGroup).toHaveBeenCalledWith(5);
    expect(service.count()).toBe(3);
  });

  it("never asks, and shows nothing, without group.activities.read", () => {
    setUp({ 5: [] });

    expect(api.getProcessingQueueForGroup).not.toHaveBeenCalled();
    expect(service.count()).toBeUndefined();
  });

  it("keeps the old count and does not throw when the fetch fails", () => {
    setUp({ 5: [Permission.GroupActivitiesRead] });
    api.getProcessingQueueForGroup.mockReturnValue(throwError(() => new Error("boom")));

    service.refresh();

    expect(service.count()).toBe(3);
  });

  it("takes the page's own result for the selected group only", () => {
    setUp({ 5: [Permission.GroupActivitiesRead] });

    service.report(5, 7);
    expect(service.count()).toBe(7);
    service.report(6, 1);
    expect(service.count()).toBe(7);
  });

  it("does not refetch on a navigation right after a fetch, nor while on the Processing page", () => {
    setUp({ 5: [Permission.GroupActivitiesRead] });
    api.getProcessingQueueForGroup.mockClear();

    events.next(new NavigationEnd(1, "/receipts/group/5", "/receipts/group/5"));
    expect(api.getProcessingQueueForGroup).not.toHaveBeenCalled();

    jest.spyOn(Date, "now").mockReturnValue(Date.now() + 60_000);
    events.next(new NavigationEnd(2, "/processing/group/5", "/processing/group/5"));
    expect(api.getProcessingQueueForGroup).not.toHaveBeenCalled();

    events.next(new NavigationEnd(3, "/receipts/group/5", "/receipts/group/5"));
    expect(api.getProcessingQueueForGroup).toHaveBeenCalledTimes(1);
    jest.restoreAllMocks();
  });
});
