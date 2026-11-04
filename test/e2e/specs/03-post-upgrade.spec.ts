import { test } from '../helpers/fixtures'
import { shoot } from '../helpers/screenshot'
import { env } from '../helpers/env'
import { openRoot, signIn, publishVideo, expectInMyVideos } from '../helpers/peertube'

const user = env('PLAYWRIGHT_DEVICE_USER')
const password = env('PLAYWRIGHT_DEVICE_PASSWORD')
const video = env('PLAYWRIGHT_VIDEO')
const preTitle = 'pre upgrade video'
const postTitle = 'post upgrade video'

test.describe('peertube post-upgrade', () => {
  test('pre-upgrade video survives and publishing still works', async ({ page }, testInfo) => {
    await openRoot(page)
    await signIn(page, user, password)
    await expectInMyVideos(page, [preTitle])
    await shoot(page, testInfo, 'pre-upgrade-video-verified')
    await publishVideo(page, video, postTitle)
    await expectInMyVideos(page, [preTitle, postTitle])
    await shoot(page, testInfo, 'post-upgrade-library')
  })
})
