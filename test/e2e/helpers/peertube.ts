import { Page, expect } from '@playwright/test'

export async function openRoot(page: Page) {
  await page.goto('/')
  await expect(page.getByRole('banner').getByRole('link', { name: 'Login' })).toBeVisible()
}

async function authelia(page: Page, user: string, password: string) {
  await page.locator('#username-textfield').fill(user)
  await page.locator('#password-textfield').fill(password)
  await page.locator('#sign-in-button').click()
}

export async function signIn(page: Page, user: string, password: string) {
  await page.getByRole('banner').getByRole('link', { name: 'Login' }).click()
  await authelia(page, user, password)
}

export async function dismissWelcomeWizard(page: Page) {
  const wizard = page.getByRole('dialog')
  await expect(wizard.getByText('Welcome to PeerTube', { exact: true })).toBeVisible()
  await wizard.getByRole('button', { name: 'Do not display again' }).click()
  await expect(wizard).toHaveCount(0)
}

function mainMenu(page: Page) {
  return page.getByRole('navigation', { name: 'Main menu' })
}

export async function publishVideo(page: Page, file: string, title: string) {
  await mainMenu(page).getByRole('link', { name: 'Publish' }).click()
  await page.locator('#videofile').setInputFiles(file)
  await expect(page.getByText('File uploaded! Your video is now available in your private library.')).toBeVisible({ timeout: 180_000 })
  await page.getByRole('textbox', { name: 'Video title' }).fill(title)
  await page.getByRole('button', { name: 'Save', exact: true }).click()
  await expect(page.getByText('Changes saved.')).toBeVisible()
}

export async function openMyVideos(page: Page) {
  await mainMenu(page).getByRole('link', { name: 'Videos', exact: true }).click()
}

export async function expectInMyVideos(page: Page, titles: string[]) {
  await openMyVideos(page)
  for (const title of titles) {
    await expect(page.getByRole('link', { name: title, exact: true }).first()).toBeVisible()
  }
}
