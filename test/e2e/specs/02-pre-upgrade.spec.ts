import { test } from '../helpers/fixtures'
import { shoot } from '../helpers/screenshot'
import { env } from '../helpers/env'
import { openRootV6, signInV6, remindWelcomeLaterV6, publishVideoV6, expectInMyVideosV6 } from '../helpers/peertube'

const user = env('PLAYWRIGHT_DEVICE_USER')
const password = env('PLAYWRIGHT_DEVICE_PASSWORD')
const video = env('PLAYWRIGHT_VIDEO')
const title = 'pre upgrade video'

test.describe('peertube pre-upgrade', () => {
  test('seed a video on the released version', async ({ page }, testInfo) => {
    await openRootV6(page)
    await signInV6(page, user, password)
    await remindWelcomeLaterV6(page)
    await shoot(page, testInfo, 'pre-upgrade-main')
    await publishVideoV6(page, video, title)
    await expectInMyVideosV6(page, [title])
    await shoot(page, testInfo, 'pre-upgrade-library')
  })
})
