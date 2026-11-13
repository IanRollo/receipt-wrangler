import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { GroupGuard } from "../guards/group.guard";
import { ProcessingQueueListComponent } from "./processing-queue-list/processing-queue-list.component";
import { processingQueueGuard } from "./processing-queue.guard";

const routes: Routes = [
  {
    path: "group/:groupId",
    component: ProcessingQueueListComponent,
    canActivate: [GroupGuard, processingQueueGuard],
    data: {
      groupGuardBasePath: `/processing/group`,
    },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProcessingQueueRoutingModule {}
