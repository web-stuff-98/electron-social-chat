<template>
  <div class="messages">
    <div class="message">
      <p>
        <b>Username: </b>
        Lorem ipsum dolor sit amet consectetur adipisicing elit. Iusto vitae
        commodi adipisci corporis, doloribus aperiam odit possimus aspernatur
        saepe eaque.
      </p>
    </div>
  </div>
</template>
<style scoped>
.messages {
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
  padding: var(--padding-medium);
  box-sizing: border-box;
  top: 0;
  left: 0;
}
.message {
  width: 100%;
  height: 100%;
  margin: 0;
}
p {
  text-align: left;
  margin: 0;
  font-size: 0.8rem;
}
b {
  filter: opacity(0.5);
  margin: 0;
  font-size: 0.8rem;
}
</style>
