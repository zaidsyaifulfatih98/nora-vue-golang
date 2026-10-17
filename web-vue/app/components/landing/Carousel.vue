<script setup lang="ts">
const props = defineProps<{
  items: unknown[]
  /** Kalau diisi, carousel hanya aktif di bawah breakpoint ini dan berubah jadi grid di atasnya. */
  gridFrom?: 'sm' | 'lg'
  trackClass?: string
}>()
const { t } = useI18n()
const trackRef = ref<HTMLDivElement | null>(null)

function scroll(direction: 1 | -1) {
  const track = trackRef.value
  if (!track) return
  const card = track.firstElementChild as HTMLElement | null
  const amount = (card?.offsetWidth ?? track.clientWidth) + 24
  track.scrollBy({ left: direction * amount, behavior: 'smooth' })
}

const showArrows = computed(() => props.items.length > 1)

// Class Tailwind harus statis agar ikut ter-generate, jadi dipetakan per breakpoint.
const GRID_CLASSES = {
  sm: {
    arrows: 'sm:hidden',
    track: 'sm:grid sm:grid-cols-2 sm:overflow-visible lg:grid-cols-3',
    item: 'sm:w-auto',
  },
  lg: {
    arrows: 'lg:hidden',
    track: 'lg:grid lg:grid-cols-3 lg:overflow-visible',
    item: 'lg:w-auto',
  },
} as const
const grid = computed(() => (props.gridFrom ? GRID_CLASSES[props.gridFrom] : null))
</script>

<template>
  <div class="relative">
    <button
      v-if="showArrows"
      :aria-label="t('landing.gallery.prev')"
      :class="grid?.arrows"
      class="absolute top-1/2 left-0 z-10 -translate-x-4 -translate-y-1/2 rounded-full bg-white p-3 text-[#1E2537] shadow-md ring-1 ring-[#E4E2DC] transition hover:bg-[#F0EFEA] sm:-translate-x-6"
      @click="scroll(-1)"
    >
      <Icon name="heroicons:chevron-left" class="text-xl" />
    </button>

    <div
      ref="trackRef"
      :class="[grid?.track, trackClass]"
      class="flex snap-x snap-mandatory gap-6 overflow-x-auto scroll-smooth pb-2 [-ms-overflow-style:none] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
    >
      <div
        v-for="(item, i) in items"
        :key="i"
        class="w-[85%] shrink-0 snap-start"
        :class="grid ? grid.item : 'sm:w-[46%] lg:w-[31.5%]'"
      >
        <slot :item="item" />
      </div>
    </div>

    <button
      v-if="showArrows"
      :aria-label="t('landing.gallery.next')"
      :class="grid?.arrows"
      class="absolute top-1/2 right-0 z-10 -translate-y-1/2 translate-x-4 rounded-full bg-white p-3 text-[#1E2537] shadow-md ring-1 ring-[#E4E2DC] transition hover:bg-[#F0EFEA] sm:translate-x-6"
      @click="scroll(1)"
    >
      <Icon name="heroicons:chevron-right" class="text-xl" />
    </button>
  </div>
</template>
