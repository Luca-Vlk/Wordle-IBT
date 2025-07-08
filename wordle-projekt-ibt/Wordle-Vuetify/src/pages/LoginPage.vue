
<template>
  <main>
    <v-container class="d-flex justify-center pa-0 mt-8">
      <v-sheet class="border rounded pa-8 ma-0" width="700" elevation="20">
        <div class="d-flex justify-center title">
          Wordle Login
        </div>
        
        <div class="d-flex justify-center combobox">
          <div style="width: 200px;">
            <v-combobox
                prepend-icon="mdi-account"
                variant="solo-filled"
                v-model:model-value="selectedName"
                label=""
                autocomplete="off"
                :items="nameArray">
              </v-combobox>
          </div>
        </div>

        <div class="d-flex justify-center mt-3 mb-3">
          <v-btn
            prepend-icon="mdi-login"
            :disabled="selectedName == ''"
            color="primary"
            @click="loginMethod"
          >
              Login button
          </v-btn>

        </div>
      </v-sheet>
    </v-container>
  </main>
</template>

<script setup>
import {onMounted, ref, watch} from "vue";
import { useRouter } from 'vue-router';

const router = useRouter();
const selectedName = ref("");
const nameArray = ref([""]);
onMounted(() => {
  db_selectNames().then(data => nameArray.value = data.split(","));
});

async function db_selectNames() {
  return fetch("http://localhost:8080/api/selectNames")
      .then(response => response.text())
      .then(data => {
        console.log("Userliste: ", data);
        return data;
      })
      .catch(error => {
        console.error("Fehler beim Abrufen:", error);
      });
}

function loginMethod() {
  console.log("Ausgewählter User: " + selectedName.value);

  router.push({
    name: 'home', 
    params:{user:selectedName.value}
  });
}
</script>

<style>
  .title{
    font-size: 6rem;
    font-weight: 300;
    line-height: 1;
    letter-spacing: -0.015625em;
    margin-top: 25px
  }
  .combobox{
    margin-top: 90px;
  }
</style>