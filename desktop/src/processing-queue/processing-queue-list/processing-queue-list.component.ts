import { AfterViewInit, Component, computed, DestroyRef, effect, inject, signal, TemplateRef, untracked, viewChild } from "@angular/core";
import { toSignal } from "@angular/core/rxjs-interop";
import { MatDialog } from "@angular/material/dialog";
import { MatTableDataSource } from "@angular/material/table";
import { ActivatedRoute } from "@angular/router";
import { Store } from "@ngxs/store";
import { filter, map, Subscription, take } from "rxjs";
import { DEFAULT_DIALOG_CONFIG, DEFAULT_HOST_CLASS } from "../../constants";
import { ProcessingQueueItem, ProcessingQueueService } from "../../open-api";
import { SnackbarService } from "../../services/snackbar.service";
import { ConfirmationDialogComponent } from "../../shared-ui/confirmation-dialog/confirmation-dialog.component";
import {
  SourceFileViewerDialogComponent,
  SourceFileViewerDialogData,
} from "../../shared-ui/source-file-viewer-dialog/source-file-viewer-dialog.component";
import { GroupState } from "../../store";
import { TableColumn } from "../../table/table-column.interface";
import { downloadFile, filenameFromContentDisposition } from "../../utils/file";
import { ProcessingQueueCountService } from "../processing-queue-count.service";
import {
  hasInFlightItems,
  isProcessingNow,
  itemDisplayName,
  PROCESSING_POLL_INTERVAL_MS,
  PROCESSING_QUEUE_EMPTY_MESSAGE,
  processingStateLabel,
  processingStateTone,
  processingTaskTypeLabel,
  retryLabel,
  retryNeedsConfirmation,
} from "../processing-queue.utils";

/**
 * The quick scans and email uploads that have not become receipts yet: waiting,
 * running, retrying after an error, or failed. Refreshes itself every few seconds
 * while anything is still moving and the tab is visible, then goes quiet.
 *
 * Every action is keyed by the item's own `queue` + `asynqTaskId` (never its
 * `systemTaskId`, which many items lack) and judged by the server in the item's
 * own group; the row only offers what the server flagged.
 */
@Component({
  selector: "app-processing-queue-list",
  templateUrl: "./processing-queue-list.component.html",
  styleUrls: ["./processing-queue-list.component.scss"],
  host: DEFAULT_HOST_CLASS,
  standalone: false,
})
export class ProcessingQueueListComponent implements AfterViewInit {
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly countService = inject(ProcessingQueueCountService);
  private readonly destroyRef = inject(DestroyRef);
  private readonly matDialog = inject(MatDialog);
  private readonly processingQueueService = inject(ProcessingQueueService);
  private readonly snackbarService = inject(SnackbarService);
  private readonly store = inject(Store);

  private readonly stateCell = viewChild.required<TemplateRef<any>>("stateCell");
  private readonly typeCell = viewChild.required<TemplateRef<any>>("typeCell");
  private readonly nameCell = viewChild.required<TemplateRef<any>>("nameCell");
  private readonly groupCell = viewChild.required<TemplateRef<any>>("groupCell");
  private readonly uploaderCell = viewChild.required<TemplateRef<any>>("uploaderCell");
  private readonly uploadedAtCell = viewChild.required<TemplateRef<any>>("uploadedAtCell");
  private readonly attemptsCell = viewChild.required<TemplateRef<any>>("attemptsCell");
  private readonly nextRetryCell = viewChild.required<TemplateRef<any>>("nextRetryCell");
  private readonly errorCell = viewChild.required<TemplateRef<any>>("errorCell");
  private readonly actionsCell = viewChild.required<TemplateRef<any>>("actionsCell");

  public readonly groupId = toSignal(
    this.activatedRoute.paramMap.pipe(map((params) => Number(params.get("groupId")))),
    { initialValue: Number(this.activatedRoute.snapshot.paramMap.get("groupId")) }
  );

  private readonly groups = this.store.selectSignal(GroupState.groups);

  public readonly group = computed(() => this.groups().find((group) => group.id === this.groupId()));

  public readonly isAllGroup = computed(() => !!this.group()?.isAllGroup);

  public readonly headerText = computed(() => {
    const name = this.group()?.name;
    return name ? `${name} Processing` : "Processing";
  });

  public readonly subtitle = computed(() =>
    this.isAllGroup()
      ? "Quick scans and email uploads still being processed or failed, across every group you can read."
      : "Quick scans and email uploads still being processed, or that failed."
  );

  public readonly items = signal<ProcessingQueueItem[]>([]);
  public readonly truncated = signal(false);
  public readonly loaded = signal(false);
  public readonly loadFailed = signal(false);

  /** Rows with a retry request in flight, keyed by `queue/asynqTaskId`. */
  private readonly retrying = signal<ReadonlySet<string>>(new Set());

  public readonly polling = computed(() => hasInFlightItems(this.items()));

  public readonly emptyMessage = PROCESSING_QUEUE_EMPTY_MESSAGE;

  public readonly dataSource = computed(() => new MatTableDataSource<ProcessingQueueItem>(this.items()));

  /** The cell templates exist only once the view does (same as the claim list). */
  private readonly viewReady = signal(false);

  public readonly columns = computed<TableColumn[]>(() => {
    if (!this.viewReady()) {
      return [];
    }

    return this.buildColumns();
  });

  private buildColumns(): TableColumn[] {
    return [
    { columnHeader: "State", matColumnDef: "state", template: this.stateCell(), sortable: false },
    { columnHeader: "Type", matColumnDef: "type", template: this.typeCell(), sortable: false },
    { columnHeader: "File", matColumnDef: "name", template: this.nameCell(), sortable: false },
    ...(this.isAllGroup()
      ? [{ columnHeader: "Group", matColumnDef: "group", template: this.groupCell(), sortable: false }]
      : []),
    { columnHeader: "From", matColumnDef: "uploader", template: this.uploaderCell(), sortable: false },
    { columnHeader: "Uploaded", matColumnDef: "uploadedAt", template: this.uploadedAtCell(), sortable: false },
    { columnHeader: "Attempts", matColumnDef: "attempts", template: this.attemptsCell(), sortable: false },
    { columnHeader: "Next retry", matColumnDef: "nextRetry", template: this.nextRetryCell(), sortable: false },
    { columnHeader: "Last error", matColumnDef: "error", template: this.errorCell(), sortable: false },
    { columnHeader: "", matColumnDef: "actions", template: this.actionsCell(), sortable: false },
    ];
  }

  public readonly displayedColumns = computed(() => this.columns().map((column) => column.matColumnDef));

  public readonly stateLabel = processingStateLabel;
  public readonly stateTone = processingStateTone;
  public readonly typeLabel = processingTaskTypeLabel;
  public readonly displayName = itemDisplayName;
  public readonly showSpinner = isProcessingNow;
  public readonly retryText = retryLabel;

  private loadSubscription?: Subscription;
  private pollTimer?: ReturnType<typeof setTimeout>;

  constructor() {
    // A new route group is a new list: start from empty rather than flashing the
    // previous group's rows.
    effect(() => {
      this.groupId();
      untracked(() => {
        this.items.set([]);
        this.truncated.set(false);
        this.loaded.set(false);
        this.load();
      });
    });

    const onVisibilityChange = () => {
      if (document.hidden) {
        this.stopPolling();
      } else {
        this.load();
      }
    };
    document.addEventListener("visibilitychange", onVisibilityChange);
    this.destroyRef.onDestroy(() => {
      document.removeEventListener("visibilitychange", onVisibilityChange);
      this.stopPolling();
      this.loadSubscription?.unsubscribe();
    });
  }

  public ngAfterViewInit(): void {
    this.viewReady.set(true);
  }

  public refresh(): void {
    this.load();
  }

  private load(): void {
    this.stopPolling();
    this.loadSubscription?.unsubscribe();

    const groupId = this.groupId();
    this.loadSubscription = this.processingQueueService
      .getProcessingQueueForGroup(groupId)
      .pipe(take(1))
      .subscribe({
        next: (list) => {
          this.items.set(list.items ?? []);
          this.truncated.set(!!list.truncated);
          this.loadFailed.set(false);
          this.loaded.set(true);
          this.countService.report(groupId, this.items().length);
          this.schedulePoll();
        },
        // Keep what is on screen and stop polling; the Refresh button retries.
        error: () => {
          this.loadFailed.set(true);
          this.loaded.set(true);
        },
      });
  }

  private schedulePoll(): void {
    if (!this.polling() || document.hidden) {
      return;
    }

    this.pollTimer = setTimeout(() => this.load(), PROCESSING_POLL_INTERVAL_MS);
  }

  private stopPolling(): void {
    if (this.pollTimer !== undefined) {
      clearTimeout(this.pollTimer);
      this.pollTimer = undefined;
    }
  }

  private key(item: ProcessingQueueItem): string {
    return `${item.queue}/${item.asynqTaskId}`;
  }

  public isRetrying(item: ProcessingQueueItem): boolean {
    return this.retrying().has(this.key(item));
  }

  public onPreview(item: ProcessingQueueItem): void {
    this.processingQueueService
      .getProcessingQueueTaskSourceFile(item.queue, item.asynqTaskId)
      .pipe(take(1))
      .subscribe((sourceFile) => {
        const data: SourceFileViewerDialogData = {
          encodedImage: sourceFile.encodedImage,
          name: sourceFile.name,
        };
        this.matDialog.open(SourceFileViewerDialogComponent, { ...DEFAULT_DIALOG_CONFIG, data });
      });
  }

  public onDownload(item: ProcessingQueueItem): void {
    // Observed as a response: the original name is only on Content-Disposition.
    this.processingQueueService
      .downloadProcessingQueueTaskSourceFile(item.queue, item.asynqTaskId, "response")
      .pipe(take(1))
      .subscribe((response) => {
        if (response.body) {
          downloadFile(
            response.body,
            filenameFromContentDisposition(response.headers.get("Content-Disposition"), item.fileName || "download")
          );
        }
      });
  }

  public onRetry(item: ProcessingQueueItem): void {
    if (!retryNeedsConfirmation(item)) {
      this.retry(item);
      return;
    }

    const dialogRef = this.matDialog.open(ConfirmationDialogComponent, DEFAULT_DIALOG_CONFIG);
    dialogRef.componentInstance.headerText = "Resubmit Upload";
    dialogRef.componentInstance.dialogContent =
      `Process ${itemDisplayName(item)} again? Its retries ran out, so check the error first: ` +
      "if the receipt was actually created, this adds a duplicate.";
    dialogRef
      .afterClosed()
      .pipe(take(1), filter(Boolean))
      .subscribe(() => this.retry(item));
  }

  private retry(item: ProcessingQueueItem): void {
    const key = this.key(item);
    this.retrying.update((keys) => new Set(keys).add(key));

    this.processingQueueService
      .retryProcessingQueueTask(item.queue, item.asynqTaskId)
      .pipe(take(1))
      .subscribe({
        next: () => {
          this.snackbarService.success("Upload queued to retry.");
          this.finishRetry(key);
          this.load();
        },
        // The interceptor reports the error; the task may have moved on, so reload.
        error: () => {
          this.finishRetry(key);
          this.load();
        },
      });
  }

  private finishRetry(key: string): void {
    this.retrying.update((keys) => {
      const next = new Set(keys);
      next.delete(key);
      return next;
    });
  }
}
