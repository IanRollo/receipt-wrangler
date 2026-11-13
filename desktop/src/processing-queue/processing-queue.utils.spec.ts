import { Group, Permission } from "../open-api";
import {
  canReadProcessingQueue,
  hasInFlightItems,
  isInFlight,
  itemDisplayName,
  processingStateLabel,
  processingStateTone,
  processingTaskTypeLabel,
  retryLabel,
  retryNeedsConfirmation,
} from "./processing-queue.utils";

describe("processing-queue.utils", () => {
  it("labels and colours the known states", () => {
    expect(processingStateLabel("QUEUED")).toBe("Queued");
    expect(processingStateTone("RETRYING")).toBe("yellow");
    expect(processingStateTone("FAILED")).toBe("red");
    expect(processingStateTone("QUEUED")).toBe("gray");
  });

  it("renders an unknown state or type as its raw code", () => {
    expect(processingStateLabel("PAUSED")).toBe("PAUSED");
    expect(processingStateTone("PAUSED")).toBe("gray");
    expect(processingTaskTypeLabel("FAX_UPLOAD")).toBe("FAX_UPLOAD");
    expect(processingTaskTypeLabel("QUICK_SCAN")).toBe("Quick scan");
  });

  it("treats everything but FAILED, including an unknown state, as in flight", () => {
    expect(isInFlight({ state: "PROCESSING" })).toBe(true);
    expect(isInFlight({ state: "PAUSED" })).toBe(true);
    expect(isInFlight({ state: "FAILED" })).toBe(false);
    expect(hasInFlightItems([{ state: "FAILED" }])).toBe(false);
    expect(hasInFlightItems([{ state: "FAILED" }, { state: "QUEUED" }])).toBe(true);
  });

  it("words the retry action by state and shows it only when the server allows it", () => {
    expect(retryLabel({ state: "RETRYING", canRetryNow: true, canBeRestarted: false })).toBe("Retry now");
    expect(retryLabel({ state: "FAILED", canRetryNow: false, canBeRestarted: true })).toBe("Resubmit");
    expect(retryLabel({ state: "FAILED", canRetryNow: false, canBeRestarted: false })).toBeUndefined();
    expect(retryNeedsConfirmation({ state: "FAILED" })).toBe(true);
    expect(retryNeedsConfirmation({ state: "RETRYING" })).toBe(false);
  });

  it("falls back from file name to email subject", () => {
    expect(itemDisplayName({ fileName: "a.png", emailSubject: "Hi" })).toBe("a.png");
    expect(itemDisplayName({ fileName: "", emailSubject: "Hi" })).toBe("Hi");
    expect(itemDisplayName({ fileName: "" })).toBe("(no file name)");
  });

  describe("canReadProcessingQueue", () => {
    const real = { id: 5, isAllGroup: false } as Group;
    const other = { id: 6, isAllGroup: false } as Group;
    const all = { id: 9, isAllGroup: true } as Group;
    const read = [Permission.GroupActivitiesRead];

    it("judges a real group by its own grant", () => {
      expect(canReadProcessingQueue([all, real], { 5: read }, 5)).toBe(true);
      expect(canReadProcessingQueue([all, real, other], { 5: read }, 6)).toBe(false);
    });

    it("offers the All group when any real group grants it, and ignores the All grant itself", () => {
      expect(canReadProcessingQueue([all, real, other], { 6: read }, 9)).toBe(true);
      expect(canReadProcessingQueue([all, real], { 9: read }, 9)).toBe(false);
    });

    it("is false without a group", () => {
      expect(canReadProcessingQueue([real], { 5: read }, NaN)).toBe(false);
    });
  });
});
