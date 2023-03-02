<template>
  <div class="message">
    <p>
      <b>Username: </b>
      Lorem ipsum dolor sit amet consectetur adipisicing elit. Iusto vitae
      commodi adipisci corporis, doloribus aperiam odit possimus aspernatur
      saepe eaque.
    </p>
  </div>
</template>
<style scoped>
.container {
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
  position: absolute;
  padding: var(--padding-medium);
  box-sizing: border-box;
}
.message {
  width: 100%;
  height: 100%;
  margin: 0;
}
p {
  text-align: left;
  margin: 0;
}
b {
  filter: opacity(0.5);
  margin:0;
}
</style>
