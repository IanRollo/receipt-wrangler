import { effect, inject, Injectable, signal, untracked } from "@angular/core";
import { NavigationEnd, Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { catchError, filter, of, Subject, switchMap } from "rxjs";
import { ProcessingQueueService } from "../open-api";
import { AuthState, GroupState } from "../store";
import { canReadProcessingQueue } from "./processing-queue.utils";

/** Navigations closer together than this reuse the count already on screen. */
export const PROCESSING_COUNT_MIN_INTERVAL_MS = 15_000;

const PROCESSING_PATH = "/processing";

/**
 * The header badge: how many uploads in the selected group are in flight or
 * failed. Every item the queue lists is one of those (finished uploads are
 * receipts, not items), so the count is the list's length.
 *
 * It deliberately does not poll. It is fetched when the selected group (or the
 * permission) changes, on navigation (throttled), and when something that knows
 * it just changed the queue says so: Quick Scan after a submit (`refresh()`),
 * and the Processing page after each of its own loads (`report()`). The page
 * itself polls while something is in flight; the badge rides on that.
 */
@Injectable({ providedIn: "root" })
export class ProcessingQueueCountService {
  private readonly processingQueueService = inject(ProcessingQueueService);
  private readonly router = inject(Router);
  private readonly store = inject(Store);

  private readonly isLoggedIn = this.store.selectSignal(AuthState.isLoggedIn);
  private readonly groups = this.store.selectSignal(GroupState.groups);
  private readonly groupPermissions = this.store.selectSignal(AuthState.groupPermissions);
  private readonly selectedGroupId = this.store.selectSignal(GroupState.selectedGroupId);

  private readonly request$ = new Subject<number>();
  private lastFetchAt = 0;

  /** Undefined until known, and for a group the viewer cannot see the queue of. */
  public readonly count = signal<number | undefined>(undefined);

  constructor() {
    this.request$
      .pipe(
        switchMap((groupId) =>
          this.processingQueueService.getProcessingQueueForGroup(groupId).pipe(
            // A badge must never be the reason a page shows an error.
            catchError(() => of(undefined))
          )
        )
      )
      .subscribe((list) => {
        if (list) {
          this.count.set(list.items.length);
        }
      });

    // A new group (or newly arrived permissions) is a new question: drop the old
    // answer and ask again.
    effect(() => {
      this.selectedGroupId();
      const visible = this.isVisible();
      untracked(() => {
        this.count.set(undefined);
        if (visible) {
          this.refresh();
        }
      });
    });

    this.router.events
      .pipe(filter((event): event is NavigationEnd => event instanceof NavigationEnd))
      .subscribe((event) => {
        const onProcessingPage = event.urlAfterRedirects.startsWith(PROCESSING_PATH);
        const fresh = Date.now() - this.lastFetchAt < PROCESSING_COUNT_MIN_INTERVAL_MS;
        if (!onProcessingPage && !fresh && this.isVisible()) {
          this.refresh();
        }
      });
  }

  /** Whether the viewer may see the selected group's queue at all. */
  public isVisible(): boolean {
    return (
      this.isLoggedIn() &&
      canReadProcessingQueue(this.groups(), this.groupPermissions(), Number(this.selectedGroupId()))
    );
  }

  /** Fetches the selected group's count now. */
  public refresh(): void {
    if (!this.isVisible()) {
      return;
    }

    this.lastFetchAt = Date.now();
    this.request$.next(Number(this.selectedGroupId()));
  }

  /** The Processing page already holds a fresh list for `groupId`: use it. */
  public report(groupId: number, itemCount: number): void {
    if (groupId === Number(this.selectedGroupId())) {
      this.lastFetchAt = Date.now();
      this.count.set(itemCount);
    }
  }
}
