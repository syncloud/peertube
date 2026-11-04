import { test } from '../helpers/fixtures'
import { shoot } from '../helpers/screenshot'
import { env } from '../helpers/env'
import { openRoot, signIn, dismissWelcomeWizard, publishVideo, expectInMyVideos } from '../helpers/peertube'

const user = env('PLAYWRIGHT_DEVICE_USER')
const password = env('PLAYWRIGHT_DEVICE_PASSWORD')
const video = env('PLAYWRIGHT_VIDEO')
const title = 'pre upgrade video'

test.describe('peertube pre-upgrade', () => {
  test('seed a video on the released version', async ({ page }, testInfo) => {
    await openRoot(page)
    await signIn(page, user, password)
    await dismissWelcomeWizard(page)
    await shoot(page, testInfo, 'pre-upgrade-main')
    await publishVideo(page, video, title)
    await expectInMyVideos(page, [title])
    await shoot(page, testInfo, 'pre-upgrade-library')
  })
})
