# Fork Roadmap

Features this fork means to explore, and what has been built.

## Built

- **Reimbursement claims** (branch `claims-tracking`, 2026-11-11). A claim bundles a group's
  receipts submitted together, with a status, submitted/due dates, a claimed amount and a dated
  ledger of payments, recoupments and notes; the server computes claimed, paid, outstanding and
  days waiting. This is item 3 below made concrete: a recoupment is the "differential" recorded
  against the claim it hit. See `CLAUDE.md` → "Reimbursement Claims".
- **Claim packet** (2026-11-12). One zip per claim: `packet.pdf` (cover sheet with the figures,
  the ledger, receipts with category subtotals, then every receipt with its images), the original
  image files untouched, `receipts.csv` / `items.csv` with an Images column, and a `manifest.json`
  carrying a SHA-256 for every file. This delivers items 1 and 2 below for claimed receipts; the
  plain receipt export is unchanged.

- **Multicurrency, slice A** (branch `multicurrency`, 2026-11-12). Backend only so far: an explicit
  ISO 4217 currency on every receipt (request, else group default, else install default), amounts
  validated by the currency's exponent, a recorded conversion (base amount, rate, date, source) kept
  as evidence, the Receipt Summary's legacy total restricted to the base currency with per-currency
  totals beside it, claims reckoned in the group's base currency with unconverted receipts set
  aside, and the claim packet and CSV showing all of it. See item 4 and `api/CLAUDE.md` →
  "Multicurrency (fork feature)".

**Why this fork exists:** a personal record of payments that can be handed over as proof,
exact transactions with their receipt images, independent of the system of whoever is disputing
them. Features are judged by how well they produce that evidence.

## 1. Include images in the receipt export (done for claims via the packet; plain export still CSV-only)

Today **Export to CSV** (filtered set or selected rows) zips `receipts.csv` + `items.csv` only;
images can only be downloaded one at a time from each receipt.

- Add an `images/` folder to the export zip, with descriptive names such as
  `2026-10-03_Costco_412_1.jpg`.
- Add an **Images** column to `receipts.csv` listing each row's filenames.
- Possibly an "include images" checkbox so a data-only export stays light.
- Add a **SHA-256** column per image to show files are unchanged since they were recorded.

Where: `api/internal/services/receipt_csv.go` (`GetZippedCsvFiles`), both handlers in
`api/internal/handlers/export.go`. Filtering and permissions are already handled there.

## 2. One-PDF evidence bundle (done, as the claim packet)

One page per transaction: date, amount, name and the image together, as a single PDF that is
easy to hand over. Reuse the HTML-to-PDF path the Report Builder already uses
(`NewHtmlToPdfService`, Chromium). Consider a cover page with totals and the date range.

## 3. Discrepancy tracking ("differentials") — covered by the claim ledger

Recording what was billed against what is actually owed, and surfacing the difference. Likely
built on custom fields plus the existing Report Builder / Receipt Summary totals rather than a new
subsystem. The receipt change history already shows when each receipt was recorded.

## 4. Multicurrency support (large, cross-cutting)

Today currency is **display-only and install-wide**: `SystemSettings.CurrencyDisplay` (symbol,
separators, position) formats every amount the same way. `Receipt.Amount`, `Item.Amount` and
CURRENCY custom field values are bare decimals with no currency attached, so mixed currencies would
silently sum as if they were one.

Likely phases:

- **Record it.** A currency code (ISO 4217) per receipt, with items inheriting the receipt's. A
  default per install or group. The AI extraction prompt returns the currency it reads, and
  displays and exports show the code next to every amount. On the wire it must be a **plain string,
  never an enum** (the closed-Dart-enum outage pattern in `CLAUDE.md`).
- **Never sum across currencies.** Receipt Summary, Report Builder, the who-owes-whom settlement
  and dashboard widgets either group totals by currency or convert. This is where most of the work
  is.
- **Convert (optional).** Store the original amount, its currency, the rate used and its date on
  the receipt rather than re-deriving them later, so a total can be explained in a dispute. That
  needs a rate source.
- **Precision.** `Receipt.Amount` is `decimal(10,2)`; JPY uses 0 decimal places and KWD/BHD use 3.

Touches all three components plus a `swagger.yml` change with both clients regenerated in the same
change. Rough size: "record it + group totals by currency" is a week or so. Conversion is extra.

**Slice A (backend, built):** the currency per receipt with the default chain, a `GET /currency`
table, exponent validation and `decimal(13,3)` storage, the AI prompt returning the currency (and
the stored default prompt upgraded when unedited), the recorded conversion with derivation and a
`conversionNote`, the backfill of existing receipts, the Receipt Summary split by currency, claims
reckoned in the base currency with an unconverted set, the packet and CSV columns, history entries.
Also the "Convert (optional)" phase above, minus a rate source: the user types the rate or the
amount charged; nothing is fetched.

**Still to do:**

- Desktop UI for all of the above (in progress), and mobile (only the regenerated client so far).
- **Slice B:** the Report Builder (a `currency` dimension, per-currency sums, 3-place split
  scale), the dashboard pie chart, and who-owes-whom (netted per currency). These still sum bare
  amounts.
- A rate source, if one is ever wanted; custom CURRENCY field values still have no currency of
  their own (they inherit the receipt's in the summary and packet).
- Re-deriving stored conversions if a group's base currency is changed.

## 5. Review queue for uncertain scans

When a quick scan or email upload can't confidently read something (an unclear currency, a smudged
total, an ambiguous date), flag the receipt with the specific concern instead of silently guessing.
The user then works through a review queue to resolve each concern.

What already exists:

- A `NEEDS_ATTENTION` receipt status, and per-group/per-user default statuses for quick scan and
  email (`QuickScanDefaultStatus`, `EmailDefaultReceiptStatus`). Setting the default to
  `NEEDS_ATTENTION` and filtering on it is a crude queue available today, with no code.
- Editable prompts (`api/internal/services/prompts.go`). The default prompt says "Omit any value if
  not found with confidence", so uncertainty is currently **dropped silently**, not reported.
- Receipt comments, which can carry a note to the reviewer.

Likely phases:

- **MVP (~1-2 days, backend only).** The prompt asks for a `concerns` array (field + question). If
  it is non-empty, the receipt is saved as `NEEDS_ATTENTION` with the concerns posted as a system
  comment. Both clients get this through the existing status filter and comments. Reuse the
  existing status rather than adding one: a new value in the status enum would break
  already-released mobile builds (closed-enum outages, see `CLAUDE.md`).
- **Deterministic checks alongside the AI's.** An LLM is a poor judge of its own uncertainty, so add
  rule-based flags: missing currency, line items not summing to the total, a future date, a likely
  duplicate (same store, date and amount).
- **Proper review UI (~1 week).** Store flags in their own table (field, question, resolved by and
  when) and show each on the receipt form with a control to answer it. Resolving the last flag
  moves the receipt out of the queue. Optionally rerun extraction with the answers added to the
  prompt (the true "ask back").

Pairs naturally with multicurrency (item 4): "currency unclear" is the first flag to build.

## 6. Processing queue for quick scans and email uploads

A quick-added receipt disappears until it is finished: there is nowhere to see that it is waiting,
being read, retrying after an error, or stuck. Wanted: a page listing every upload still in flight
or failed, with its state, its image, the error if any, and a resubmit button.

What already exists:

- **Nothing is recorded while a scan runs.** `wranglerasynq/quick_scan_handler.go` calls
  `ReceiptService.QuickScan`, and the `system_tasks` rows (QUICK_SCAN, OCR_PROCESSING,
  CHAT_COMPLETION) are written as each step **finishes**. `SystemTaskStatus` is only `SUCCEEDED` /
  `FAILED`, so "queued" and "in progress" have no row to show.
- **The in-flight state lives in Redis.** Asynq already knows every task's state (pending, active,
  retry with its next attempt time and last error, archived), and the code already uses
  `asynq.Inspector` (`wranglerasynq/inspector.go`, `source_file.go`).
- **Failures are already recoverable.** "Failed Uploads Keep Their Image" (root `CLAUDE.md`): a failed
  activity keeps its source file, offers preview/download, and an archived task can be rerun
  (`canBeRestarted`). Today that is reached through the group activity feed and System Tasks, mixed
  in with everything else.

Likely shape:

- **Read the queue from asynq, not from new DB rows.** A new endpoint lists the quick-scan and email
  tasks in the pending, active, retry and archived sets, decodes each payload for its group,
  uploader and original file name, and filters to what the caller may see (the task's group, as the
  activity feed already judges it). Each row: state, uploaded at, file name, attempt count, next
  retry, last error, and links to preview the image and to rerun.
- **Do not add `QUEUED` / `IN_PROGRESS` to `SystemTaskStatus`.** It is a closed enum on the wire;
  a new value is the released-Android-build outage described in `CLAUDE.md`. The queue's own state
  should be a **plain string** on a new schema instead.
- **Resubmit reuses the existing rerun**, which already enforces `Archived` and needs the source
  file. A task still in the retry set could get a "retry now" (`Inspector.RunTask`) as well.
- **Desktop:** a "Processing" page (or a panel on the dashboard) that polls every few seconds while
  anything is in flight, plus a count badge in the header. Mobile could follow, since the person
  quick-adding is often on the phone.
- **Limit to note:** succeeded quick-scan tasks leave Redis immediately (no `asynq.Retention` on
  that queue), so the queue shows only unfinished and failed work. Finished receipts are already in
  the receipts list; this page is for what is not.

Rough size: backend endpoint and desktop page in 2-3 days; mobile extra.

**Status: backend built (2026-11-12); the desktop page is next.** `GET /processingQueue/group/{groupId}`
lists the quick-scan and email tasks in asynq's pending, scheduled, active, retry and archived sets, with
a plain-string state (`QUEUED` / `PROCESSING` / `RETRYING` / `FAILED`), attempt counts, next retry, last
error and whether the file and a rerun are available. It is gated by the activity feed's own permissions
in each item's group. Preview, download and "retry now" are keyed by queue and asynq task id, not system
task id. See `api/CLAUDE.md` → "Processing queue (fork feature)". Corrections to the plan above:

- **The existing preview/download/rerun routes are not enough.** They address a `system_tasks` row, and
  a pending or active task has none, nor does one that fails before any AI step (for example no
  receipt processing settings configured). Hence the `/processingQueue/task/{queue}/{taskId}/...` routes;
  `systemTaskId` is on an item only when a row exists.
- **"Uploaded at" is not recorded by asynq.** It is the stored file's modification time.
- **A succeeded quick scan leaves Redis at once, but a succeeded email task lingers 6 hours**
  (`completedTaskRetention`); the queue deliberately skips completed tasks, so neither shows.
- **Retry now fits cleanly** (`RETRYING` only moves the attempt forward; `FAILED` follows the rerun
  endpoint's rules), so it was built rather than left out.
- At most 100 tasks are read per queue state, and a response flag `truncated` says when that bit.

## Considered and deferred: Google Drive mirror

Nightly copy into `Receipts/<yyyy-mm>/` with a monthly spreadsheet linking to each image.
Moderate effort (OAuth with the `drive.file` scope and a published app, a sync job, a table of
what's been uploaded). **Deleting local copies after upload is out of scope**: images are read
from local disk throughout the API, so it would need a storage-backend rewrite. Revisit only if a
backup offsite becomes the priority; `rclone` on the VM is the no-code stopgap.
