import { Component, computed, effect, inject, signal, untracked } from "@angular/core";
import { Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { take, tap } from "rxjs";
import { LayoutState } from "src/store/layout.state";
import { ToggleIsSidebarOpen } from "src/store/layout.state.actions";
import { AuthService, NotificationsService, Permission } from "../../open-api";
import { canReadProcessingQueue } from "../../processing-queue/processing-queue.utils";
import { ProcessingQueueCountService } from "../../processing-queue/processing-queue-count.service";
import { AuthState, GroupState } from "../../store";

@Component({
    selector: "app-header",
    templateUrl: "./header.component.html",
    styleUrls: ["./header.component.scss"],
    standalone: false
})
export class HeaderComponent {
  public isLoggedIn = this.store.selectSignal(AuthState.isLoggedIn);

  public selectedGroupId = this.store.selectSignal(GroupState.selectedGroupId);

  public selectedGroupIdNumber = computed(() => Number(this.selectedGroupId()));

  public loggedInUser = this.store.selectSignal(AuthState.loggedInUser);

  public showProgressBar = this.store.selectSignal(LayoutState.showProgressBar);

  public userPreferences = this.store.selectSignal(AuthState.userPreferences);

  public receiptHeaderLink = computed(() => {
    this.selectedGroupId();
    return [this.store.selectSnapshot(GroupState.receiptListLink)];
  });

  public dashboardHeaderLink = computed(() => {
    this.selectedGroupId();
    return [this.store.selectSnapshot(GroupState.dashboardLink)];
  });

  // Reimbursement claims (fork feature). Built here rather than as a GroupState
  // selector beside receiptListLink, to keep the fork's footprint in shared
  // upstream files to this one line of logic.
  public claimsHeaderLink = computed(() => ["/claims/group", this.selectedGroupId()]);

  // Processing queue (fork feature). The count service owns the badge and its
  // refresh rules; the link shows only where the viewer may read the queue.
  private readonly processingQueueCount = inject(ProcessingQueueCountService);

  public processingHeaderLink = computed(() => ["/processing/group", this.selectedGroupId()]);

  public canViewProcessing = computed(() =>
    canReadProcessingQueue(this.groups(), this.groupPermissions(), this.selectedGroupIdNumber())
  );

  public processingCount = this.processingQueueCount.count;

  private readonly groups = this.store.selectSignal(GroupState.groups);

  private readonly groupPermissions = this.store.selectSignal(AuthState.groupPermissions);

  public settingsBaseHeaderLink = computed(() => {
    this.selectedGroupId();
    return [this.store.selectSnapshot(GroupState.settingsLinkBase) + "/view"];
  });

  public groupName = computed(() => {
    const groupId = this.selectedGroupId();
    const group = this.store.selectSnapshot(GroupState.getGroupById(groupId));
    return group?.name as string ?? "";
  });

  public notificationCount = signal<number | undefined>(undefined);

  protected readonly Permission = Permission;

  private readonly canReadNotifications = this.store.selectSignal(
    AuthState.hasAppPermission(Permission.AppNotificationsRead)
  );

  constructor(
    private authService: AuthService,
    private notificationsService: NotificationsService,
    private router: Router,
    private store: Store
  ) {
    this.listenForLoggedInUser();
  }

  private listenForLoggedInUser(): void {
    let fetched = false;
    effect(() => {
      const loggedIn = this.isLoggedIn();
      const canRead = this.canReadNotifications();
      if (!loggedIn) {
        fetched = false;
        return;
      }
      if (canRead && !fetched) {
        fetched = true;
        untracked(() => {
          this.notificationsService.getNotificationCount().pipe(
            take(1),
            tap((n) => {
              this.notificationCount.set(n > 0 ? n : undefined);
            })
          ).subscribe();
        });
      }
    });
  }

  public toggleSidebar(): void {
    this.store.dispatch(new ToggleIsSidebarOpen());
  }
}
