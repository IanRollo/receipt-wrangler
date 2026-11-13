import { inject } from "@angular/core";
import { CanActivateFn, Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { AuthState, GroupState } from "../store";
import { canReadProcessingQueue } from "./processing-queue.utils";

/**
 * Runs after `GroupGuard` has settled the route group. A real group needs
 * `group.activities.read` itself; the All group needs it in any real group.
 * Otherwise the user lands on their dashboard, as for every other group gate.
 */
export const processingQueueGuard: CanActivateFn = (route) => {
  const store = inject(Store);
  const router = inject(Router);

  const allowed = canReadProcessingQueue(
    store.selectSnapshot(GroupState.groups),
    store.selectSnapshot(AuthState.groupPermissions),
    Number(route.params["groupId"])
  );

  return allowed || router.parseUrl(store.selectSnapshot(GroupState.dashboardLink));
};
