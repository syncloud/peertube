import { Page, expect } from '@playwright/test'

export async function openRoot(page: Page) {
  await page.goto('/')
  await expect(page.getByRole('banner').getByRole('link', { name: 'Login' })).toBeVisible()
}

export async function openRootV6(page: Page) {
  await page.goto('/')
  await expect(mainMenu(page).getByRole('link', { name: 'Login' })).toBeVisible()
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

export async function signInV6(page: Page, user: string, password: string) {
  await mainMenu(page).getByRole('link', { name: 'Login' }).click()
  await authelia(page, user, password)
}

export async function dismissWelcomeWizard(page: Page) {
  const wizard = page.getByRole('dialog')
  await expect(wizard.getByText('Welcome to PeerTube', { exact: true })).toBeVisible()
  await wizard.getByRole('button', { name: 'Do not display again' }).click()
  await expect(wizard).toHaveCount(0)
}

export async function dismissCollaborationAnnouncement(page: Page) {
  const modal = page.getByRole('dialog')
  await expect(modal.getByRole('heading', { name: 'Collaboration on channels are coming to PeerTube' })).toBeVisible()
  await modal.getByRole('button', { name: 'Got it' }).click()
  await expect(modal).toHaveCount(0)
}

export async function remindWelcomeLaterV6(page: Page) {
  const modal = page.getByRole('dialog')
  await expect(modal.getByRole('heading', { name: 'Welcome to PeerTube, dear administrator!' })).toBeVisible()
  await modal.getByRole('button', { name: 'Remind me later' }).click()
  await expect(modal).toHaveCount(0)
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

export async function publishVideoV6(page: Page, file: string, title: string) {
  await page.getByRole('link', { name: 'Publish', exact: true }).click()
  await page.locator('#videofile').setInputFiles(file)
  await expect(page.getByText('Congratulations! Your video is now available in your private library.')).toBeVisible({ timeout: 180_000 })
  await page.locator('#name').fill(title)
  await page.getByRole('button', { name: 'Publish', exact: true }).click()
  await expect(page.getByText('Video published.')).toBeVisible()
  await expect(page.getByRole('heading', { name: title })).toBeVisible()
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

export async function expectInMyVideosV6(page: Page, titles: string[]) {
  await mainMenu(page).getByRole('link', { name: 'Videos', exact: true }).click()
  for (const title of titles) {
    await expect(page.getByRole('link', { name: `Watch video ${title}`, exact: true }).first()).toBeVisible()
  }
}
