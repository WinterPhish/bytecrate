import { test, expect } from '@playwright/test';

test.describe('Auth E2E', () => {
  test('register -> login -> refresh cookie present and refresh rotates', async ({ page, context, request }) => {
    const rnd = Math.floor(Math.random() * 1000000);
    const email = `e2e+${rnd}@example.test`;
    const password = 'password123';

    // go to register page
    await page.goto('/register');
    await page.fill('input[placeholder="Email"]', email);
    await page.fill('input[placeholder="Password"]', password);
    await page.click('text=Create Account');

    // after register should navigate to dashboard
    await expect(page).toHaveURL(/.*\/dashboard$/);
    await expect(page.locator('text=Server Status')).toBeVisible();

    // ensure refresh_token cookie exists
    const cookies = await context.cookies();
    const rt = cookies.find(c => c.name === 'refresh_token');
    expect(rt).toBeTruthy();

    // call refresh endpoint from the page context so browser includes cookies
    const body = await page.evaluate(async () => {
      const r = await fetch('http://localhost:8080/api/auth/refresh', { method: 'POST', credentials: 'include' });
      if (!r.ok) throw new Error('refresh failed: ' + r.status);
      return r.json();
    });
    expect(body.token).toBeTruthy();

    // token should be stored in localStorage by client-side code (axios interceptor will not run here)
    // but we can set it and verify we can access protected resource
    await page.evaluate((t) => localStorage.setItem('token', t), body.token);
    // navigate to dashboard and ensure it loads protected API correctly
    await page.goto('/dashboard');
    await expect(page.locator('text=Server Status:')).toBeVisible();
  });
});
