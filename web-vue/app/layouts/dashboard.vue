<template>
  <div class="grid grid-cols-12">
    <div class="col-span-3">
      <SideBar />
    </div>
    <div class="col-span-9">
      <DashboardHeader />
      <div class="p-10">
        <slot />
      </div>
    </div>
  </div>
</template>
