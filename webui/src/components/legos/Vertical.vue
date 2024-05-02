<template>
  <div class="c-vertical">
    <div class="c-head">
      <slot name="header"></slot>
    </div>

    <div class="c-main">
      <slot></slot>
    </div>

    <div class="c-footer">
      <slot name="footer"></slot>
    </div>

  </div>
</template>

<style>
.c-vertical {
  display: grid;
  grid-template-areas:
      "header"
      "content"
      "footer";;
  grid-template-rows: auto 1fr auto;
  grid-template-columns: 1fr;
  min-height: 100vh;
  width:100%;
}

.c-vertical .c-head {
  grid-area: header;
}
.c-vertical .c-main {
  grid-area: content;
}
.c-vertical .c-footer {
  grid-area: footer;
}

</style>
