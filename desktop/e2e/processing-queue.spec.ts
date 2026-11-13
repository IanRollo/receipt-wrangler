import { expect, test } from '@playwright/test';
import { readFileSync } from 'node:fs';
import { creds } from './helpers/auth';
import {
  apiCreateGroup,
  apiCreateUnreachableProcessingSettings,
  apiDeleteGroupById,
  apiDeleteProcessingSettings,
  apiGetUserId,
  apiPatchSystemSettings,
  apiQuickScan,
  uniqueName,
  withAdminApi,
  withApiAs,
} from './helpers/provisioning';

const FIXTURE = 'e2e/fixtures/receipt.png';

test.use({ storageState: 'e2e/.auth/admin.json' });

interface QueueItem {
  asynqTaskId: string;
  queue: string;
  state: string;
  retryCount: number;
  lastError: string;
  canRetryNow: boolean;
  hasSourceFile: boolean;
}

// A quick scan has to fail on purpose, and the only lever is the AI provider,
// which is a GLOBAL system setting. So this suite mutates shared server state,
// runs serially, and restores what it found, leaning on the `e2e-shared-backend`
// job-level concurrency group exactly like failed-activity-source-file.spec.ts.
//
// The failing scan is not given time to reach FAILED (asynq's backoff makes that
// minutes long); RETRYING is the state a real user sees first, and "Retry now"
// is the action that matters there.
test.describe.serial('Processing queue', () => {
  let groupId: number;
  let processingSettingsId: number;
  let promptId: number;
  let originalProcessingSettingsId: unknown;

  const readQueue = async (api: Parameters<Parameters<typeof withAdminApi>[0]>[0]) => {
    const res = await api.get(`/api/processingQueue/group/${groupId}`);
    expect(res.ok()).toBe(true);
    return (await res.json()) as { items: QueueItem[]; truncated: boolean };
  };

  test.beforeAll(async () => {
    test.setTimeout(180_000);

    await withAdminApi(async (api) => {
      const settings = await (await api.get('/api/systemSettings')).json();
      originalProcessingSettingsId = settings.receiptProcessingSettingsId;

      const adminUserId = await apiGetUserId(api, creds('admin').username);
      const created = await apiCreateUnreachableProcessingSettings(api, uniqueName('queue-scan'));
      processingSettingsId = created.id;
      promptId = created.promptId;
      await apiPatchSystemSettings(api, { receiptProcessingSettingsId: processingSettingsId });

      groupId = (await apiCreateGroup(api, uniqueName('queue-group'))).id;
      await apiQuickScan(api, groupId, adminUserId, {
        name: 'receipt.png',
        mimeType: 'image/png',
        buffer: readFileSync(FIXTURE),
      });
    });
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        await apiPatchSystemSettings(api, {
          receiptProcessingSettingsId: originalProcessingSettingsId ?? null,
        });
        if (groupId) {
          await apiDeleteGroupById(api, String(groupId));
        }
        if (processingSettingsId) {
          await apiDeleteProcessingSettings(api, processingSettingsId, promptId);
        }
      });
    } catch (error) {
      console.warn('Failed to tear down the processing-queue fixtures', error);
    }
  });

  test('a failing quick scan is listed as retrying, previews, and Retry now runs another attempt', async ({ page }) => {
    test.setTimeout(120_000);

    // Wait over the wire for the first failed attempt to be rescheduled.
    await expect
      .poll(
        async () =>
          withAdminApi(async (api) =>
            (await readQueue(api)).items.find((item) => item.state === 'RETRYING')?.retryCount ?? 0,
          ),
        { timeout: 90_000 },
      )
      .toBeGreaterThan(0);

    await page.goto(`/processing/group/${groupId}`);

    const row = page.locator('mat-row, tr.mat-mdc-row').filter({ hasText: 'receipt.png' });
    await expect(row).toHaveCount(1);
    await expect(row.getByTestId('processing-queue-state')).toContainText('Retrying');
    await expect(row.getByTestId('processing-queue-last-error')).not.toBeEmpty();
    await expect(row.getByTestId('processing-queue-attempts')).toContainText('/');
    await expect(row.getByTestId('processing-queue-preview')).toBeVisible();
    await expect(row.getByTestId('processing-queue-download')).toBeVisible();

    // The header link carries the count.
    await expect(page.getByTestId('header-processing-link')).toBeVisible();
    await expect(page.locator('[data-testid="header-processing-link"] .mat-badge-content')).toHaveText('1');

    // Same session: a second browser context from the saved state is logged out.
    const before = await withAdminApi(async (api) =>
      (await readQueue(api)).items[0],
    );
    expect(before).toBeTruthy();

    await row.getByTestId('processing-queue-preview').click();
    await expect(page.locator('app-source-file-viewer-dialog')).toBeVisible();
    await page.keyboard.press('Escape');
    await expect(page.locator('app-source-file-viewer-dialog')).toHaveCount(0);

    await row.getByTestId('processing-queue-retry').click();

    // The retry runs the attempt now and, the provider being unreachable, it
    // fails again: the attempt counter moves past what it was.
    await expect
      .poll(
        async () =>
          withAdminApi(async (api) => {
            const item = (await readQueue(api)).items.find((i) => i.asynqTaskId === before!.asynqTaskId);
            return item?.retryCount ?? 0;
          }),
        { timeout: 90_000 },
      )
      .toBeGreaterThan(before!.retryCount);
  });

  test('a user without access to the group is refused the queue', async () => {
    await withAdminApi(async (api) => {
      expect((await api.get(`/api/processingQueue/group/${groupId}`)).ok()).toBe(true);
    });
    await withApiAs('user', async (api) => {
      const res = await api.get(`/api/processingQueue/group/${groupId}`);
      expect(res.status()).toBe(403);
    });
  });
});
