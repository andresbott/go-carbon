<template>
  <div class="c-horizontal">
    <div class="c-left">
      <slot name="left"></slot>
    </div>

    <div class="c-main">
      <slot></slot>
    </div>

    <div class="c-right">
      <slot name="right"></slot>
    </div>

  </div>
</template>

<style>
.c-horizontal {
  display: grid;
  grid-template-areas: "left" "center" "right";

  grid-template-columns: auto 1fr auto;
  grid-template-rows: auto;

  width:100%;
}

.c-vertical .c-left {
  grid-area: left;
}
.c-vertical .c-main {
  grid-area: center;
}
.c-vertical .c-right {
  grid-area: right;
}

</style>
