import { Group, Permission, ProcessingQueueItem } from "../open-api";
import { hasAll } from "../utils/permission.utils";

/*
 * Processing queue (fork feature): the quick scans and email uploads that have
 * not succeeded yet. `state` and `taskType` are plain strings on the wire, so
 * every helper here tolerates a code this client has never seen: it reads as in
 * progress (state) or renders as its raw code (both), and never throws.
 */

export type ProcessingStateTone = "red" | "green" | "gray" | "yellow";

export const PROCESSING_STATE_QUEUED = "QUEUED";
export const PROCESSING_STATE_PROCESSING = "PROCESSING";
export const PROCESSING_STATE_RETRYING = "RETRYING";
export const PROCESSING_STATE_FAILED = "FAILED";

export const PROCESSING_POLL_INTERVAL_MS = 5000;

export const PROCESSING_QUEUE_EMPTY_MESSAGE =
  "Nothing is processing. Finished uploads appear in the receipt list.";

interface StateDisplay {
  label: string;
  tone: ProcessingStateTone;
}

const STATE_DISPLAY: Record<string, StateDisplay> = {
  [PROCESSING_STATE_QUEUED]: { label: "Queued", tone: "gray" },
  [PROCESSING_STATE_PROCESSING]: { label: "Processing", tone: "gray" },
  [PROCESSING_STATE_RETRYING]: { label: "Retrying", tone: "yellow" },
  [PROCESSING_STATE_FAILED]: { label: "Failed", tone: "red" },
};

const TASK_TYPE_LABELS: Record<string, string> = {
  QUICK_SCAN: "Quick scan",
  EMAIL_UPLOAD: "Email upload",
};

export function processingStateLabel(state: string | null | undefined): string {
  return STATE_DISPLAY[state ?? ""]?.label ?? state ?? "";
}

export function processingStateTone(state: string | null | undefined): ProcessingStateTone {
  return STATE_DISPLAY[state ?? ""]?.tone ?? "gray";
}

export function processingTaskTypeLabel(taskType: string | null | undefined): string {
  return TASK_TYPE_LABELS[taskType ?? ""] ?? taskType ?? "";
}

/** Only FAILED is finished business; any other state, known or not, is still moving. */
export function isInFlight(item: Pick<ProcessingQueueItem, "state">): boolean {
  return item.state !== PROCESSING_STATE_FAILED;
}

export function hasInFlightItems(items: readonly Pick<ProcessingQueueItem, "state">[]): boolean {
  return items.some(isInFlight);
}

/** Shows the spinner: actively running, or a state this client does not know. */
export function isProcessingNow(item: Pick<ProcessingQueueItem, "state">): boolean {
  return item.state === PROCESSING_STATE_PROCESSING;
}

/** What the row's retry button says, or undefined when the row offers none. */
export function retryLabel(
  item: Pick<ProcessingQueueItem, "state" | "canRetryNow" | "canBeRestarted">
): string | undefined {
  if (!item.canRetryNow && !item.canBeRestarted) {
    return undefined;
  }

  return item.state === PROCESSING_STATE_FAILED ? "Resubmit" : "Retry now";
}

/** Only a failed item is resubmitted from scratch, so only it asks first. */
export function retryNeedsConfirmation(item: Pick<ProcessingQueueItem, "state">): boolean {
  return item.state === PROCESSING_STATE_FAILED;
}

/** The file's name, else the email's subject (a body-only email has no file). */
export function itemDisplayName(
  item: Pick<ProcessingQueueItem, "fileName" | "emailSubject">
): string {
  return item.fileName || item.emailSubject || "(no file name)";
}

/**
 * Whether the page and header link are offered for a group. A real group is
 * judged by its own grant; the synthetic All group has none of its own (the
 * server lists every group where the caller holds the permission), so it is
 * offered when any real group grants it.
 */
export function canReadProcessingQueue(
  groups: readonly Group[],
  groupPermissions: { [groupId: number]: string[] },
  groupId: number | null | undefined
): boolean {
  if (!groupId) {
    return false;
  }

  const holds = (id: number) => hasAll(groupPermissions[id] ?? [], Permission.GroupActivitiesRead);
  const group = groups.find((candidate) => candidate.id === groupId);

  if (group?.isAllGroup) {
    return groups.some((candidate) => !candidate.isAllGroup && holds(candidate.id));
  }

  return holds(groupId);
}
