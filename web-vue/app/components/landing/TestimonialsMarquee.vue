<script setup lang="ts">
interface ReviewData {
  id: string
  name: string
  eventLabel: string
  quote: string
  rating: number
}

const props = defineProps<{ reviews: ReviewData[] }>()

const GRADIENTS = ['from-[#F3D9C0] to-[#E7B98F]', 'from-[#F3D9E0] to-[#E3AFC2]', 'from-[#920f0f] to-[#8F6F3E]']

function getInitials(name: string) {
  return name
    .split(/\s*&\s*|\s+/)
    .map((part) => part[0])
    .join('')
    .slice(0, 2)
    .toUpperCase()
}

const track = computed(() => [...props.reviews, ...props.reviews])
</script>

<template>
  <div class="group mt-16 overflow-hidden [mask-image:linear-gradient(to_right,transparent,black_5%,black_95%,transparent)]">
    <div class="flex w-max gap-6 animate-[marquee_40s_linear_infinite] group-hover:[animation-play-state:paused]">
      <div
        v-for="(review, i) in track"
        :key="`${review.id}-${i}`"
        class="flex w-[320px] shrink-0 flex-col rounded-2xl bg-white p-8 shadow-sm ring-1 ring-[#E4E2DC]/70 sm:w-[380px]"
      >
        <Icon name="fa6-solid:quote-left" class="text-2xl text-[#920f0f]" />
        <p class="mt-4 flex-1 font-poppins text-sm leading-relaxed text-[#39445B] italic">&ldquo;{{ review.quote }}&rdquo;</p>

        <div class="mt-6 flex items-center gap-1 text-[#920f0f]">
          <Icon v-for="r in review.rating" :key="r" name="fa6-solid:star" class="text-sm" />
        </div>

        <div class="mt-4 flex items-center gap-3 border-t border-[#F0EFEA] pt-4">
          <div class="flex h-11 w-11 items-center justify-center rounded-full bg-gradient-to-br text-sm font-bold text-white" :class="GRADIENTS[i % GRADIENTS.length]">
            {{ getInitials(review.name) }}
          </div>
          <div>
            <p class="font-poppins text-sm font-bold text-[#1E2537]">{{ review.name }}</p>
            <p class="font-poppins text-xs text-[#6C7686]">{{ review.eventLabel }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
