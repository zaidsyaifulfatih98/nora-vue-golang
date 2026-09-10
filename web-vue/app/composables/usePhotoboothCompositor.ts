import type { FrameSlot } from '~/composables/api/photoboothFrames'

function loadImage(src: string): Promise<HTMLImageElement> {
  return new Promise((resolve, reject) => {
    const img = new Image()
    img.crossOrigin = 'anonymous'
    img.onload = () => resolve(img)
    img.onerror = reject
    img.src = src
  })
}

function drawCover(ctx: CanvasRenderingContext2D, img: HTMLImageElement, x: number, y: number, w: number, h: number) {
  const imgRatio = img.width / img.height
  const boxRatio = w / h
  let sx = 0
  let sy = 0
  let sw = img.width
  let sh = img.height

  if (imgRatio > boxRatio) {
    sw = img.height * boxRatio
    sx = (img.width - sw) / 2
  } else {
    sh = img.width / boxRatio
    sy = (img.height - sh) / 2
  }

  ctx.drawImage(img, sx, sy, sw, sh, x, y, w, h)
}

// Pure canvas compositing: draws the guest's captured photos into a frame
// PNG's transparent cutout boxes, returning the final image as a data URL.
// No component state — usable from anywhere that has a frame + photos.
export function usePhotoboothCompositor() {
  async function compositeFrame(
    frameImageUrl: string,
    photoDataUrls: string[],
    slots: FrameSlot[],
    photoCount: number,
    maxCanvasSide: number,
  ): Promise<string> {
    const frameImg = await loadImage(frameImageUrl)

    // The canvas is sized to the frame PNG's own aspect ratio so its
    // transparent cutout boxes line up exactly where they were designed,
    // instead of being cropped or stretched to fit an unrelated size.
    const scale = Math.min(1, maxCanvasSide / Math.max(frameImg.width, frameImg.height))
    const canvas = document.createElement('canvas')
    canvas.width = Math.round(frameImg.width * scale)
    canvas.height = Math.round(frameImg.height * scale)
    const ctx = canvas.getContext('2d')
    if (!ctx) throw new Error('Canvas 2D context unavailable')

    ctx.fillStyle = '#1a1a1a'
    ctx.fillRect(0, 0, canvas.width, canvas.height)

    // Slots are the exact transparent cutout boxes marked in the dashboard
    // (fractions of the frame image's own dimensions). Fall back to equal
    // vertical bands only if a frame was never configured with slots.
    const effectiveSlots = slots.length
      ? slots
      : Array.from({ length: photoCount }, (_, i) => ({
          x: 0.045,
          y: 0.045 + i * (1 / photoCount),
          width: 0.91,
          height: 1 / photoCount - 0.03,
        }))

    const photoImages = await Promise.all(photoDataUrls.map((src) => loadImage(src)))
    photoImages.forEach((img, index) => {
      const slot = effectiveSlots[index]
      if (!slot) return
      drawCover(ctx, img, slot.x * canvas.width, slot.y * canvas.height, slot.width * canvas.width, slot.height * canvas.height)
    })

    ctx.drawImage(frameImg, 0, 0, canvas.width, canvas.height)
    return canvas.toDataURL('image/png')
  }

  return { compositeFrame }
}
