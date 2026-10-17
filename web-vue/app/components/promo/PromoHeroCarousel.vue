<script setup lang="ts">
interface CarouselPhoto {
  id: string
  url: string
  caption: string | null
}

const props = defineProps<{ photos: CarouselPhoto[] | null }>()

const PLACEHOLDER_SLIDES = [
  { id: 'placeholder-1', gradient: 'from-[#F3D9C0] to-[#E7B98F]', icon: 'fa6-solid:heart', caption: 'Sarah & Bima' },
  { id: 'placeholder-2', gradient: 'from-[#F3D9E0] to-[#E3AFC2]', icon: 'fa6-solid:camera', caption: 'Studio Look' },
  { id: 'placeholder-3', gradient: 'from-[#920f0f] to-[#8F6F3E]', icon: 'fa6-solid:star', caption: 'Golden Moment' },
]

const hasPhotos = computed(() => Boolean(props.photos && props.photos.length > 0))
const slideCount = computed(() => (hasPhotos.value ? props.photos!.length : PLACEHOLDER_SLIDES.length))
const index = ref(0)

let timer: ReturnType<typeof setInterval> | undefined

function startTimer() {
  if (timer) clearInterval(timer)
  if (slideCount.value <= 1) return
  timer = setInterval(() => {
    index.value = (index.value + 1) % slideCount.value
  }, 4000)
}

onMounted(startTimer)
watch(slideCount, startTimer)
onBeforeUnmount(() => timer && clearInterval(timer))
</script>

<template>
  <div class="relative aspect-[2/3] w-full max-w-[400px] overflow-hidden rounded-xl bg-white shadow-xl ring-1 ring-[#E4E2DC]">
    <template v-if="hasPhotos">
      <img
        v-for="(photo, i) in photos"
        :key="photo.id"
        :src="photo.url"
        :alt="photo.caption ?? 'Nora Photobooth'"
        class="absolute inset-0 h-full w-full object-contain transition-opacity duration-700"
        :class="i === index ? 'opacity-100' : 'opacity-0'"
      />
    </template>
    <template v-else>
      <div
        v-for="(slide, i) in PLACEHOLDER_SLIDES"
        :key="slide.id"
        class="absolute inset-0 flex flex-col items-center justify-center gap-3 bg-gradient-to-br transition-opacity duration-700"
        :class="[slide.gradient, i === index ? 'opacity-100' : 'opacity-0']"
      >
        <Icon :name="slide.icon" class="text-5xl text-white/90" />
        <p class="font-poppins text-sm font-semibold text-white">{{ slide.caption }}</p>
      </div>
    </template>

    <div v-if="slideCount > 1" class="absolute inset-x-0 bottom-0 flex justify-center gap-2 bg-gradient-to-t from-black/30 to-transparent p-4">
      <button
        v-for="i in slideCount"
        :key="i"
        type="button"
        :aria-label="`Slide ${i}`"
        class="h-2 rounded-full transition-all"
        :class="i - 1 === index ? 'w-6 bg-white' : 'w-2 bg-white/60'"
        @click="index = i - 1"
      />
    </div>
  </div>
</template>
