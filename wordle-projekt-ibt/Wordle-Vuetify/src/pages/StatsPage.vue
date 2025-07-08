<template>

  <v-app-bar app>
    <v-toolbar-title> Statistiken </v-toolbar-title>
    <v-btn icon class="mr-4" density="compact" @click="{activeTheme == 'light' ? activeTheme = 'dark' : activeTheme ='light'; theme.global.name.value = activeTheme}">
        <v-icon v-if="activeTheme == 'light'"> mdi-white-balance-sunny </v-icon>
        <v-icon v-if="activeTheme == 'dark'"> mdi-weather-night </v-icon>
    </v-btn>
    <p class="mr-4" disabled>|</p>
    <v-btn class="mr-4" density="compact" icon="mdi-home" @click="routeBacktoMain"> </v-btn>
    <v-btn class="mr-4" density="compact" icon="mdi-account" id="menu-activator">
        </v-btn>
        <v-menu activator="#menu-activator" :close-on-content-click="false">
            <v-card class="pa-1 border">
                <v-list>
                    <v-list-item-title class="ml-3"> {{ userName }} </v-list-item-title>
                </v-list>
                <v-divider></v-divider>
                <v-btn prepend-icon="mdi-logout" variant="text" @click="navigateLoginPage"> Abmelden </v-btn>
            </v-card>
        </v-menu>
  </v-app-bar>

  <v-main>
    <v-container class="d-flex justify-center pa-0">
      <v-sheet class="border pa-8 ma-0" width="700" elevation="20">
        <v-container class="autofit-width  rounded mb-8 mt-0" color="theme_color_1" elevation="10">
          <v-row v-for="m in 1">
            <v-col v-for="n in 4">
              <v-sheet
                  rounded="lg"
                  color="theme_color_1"
                  class="d-flex align-center justify-center ml-2 mr-2"
                  height="50px"
                  width="100px">
                  <v-tooltip activator="parent" location="top">{{ getToolTip(`${m}${n}`) }}</v-tooltip>
                {{ getStat(`${m}${n}`) }}
              </v-sheet>
            </v-col>
          </v-row>
        </v-container>
        
        <v-container  width="570px" class="rounded border mt-8 mb-8">
          Versuche
          <BarChart :data="barChartData" :options="barChartOptions" />
        </v-container>
      </v-sheet>
    </v-container>
  </v-main>
</template>

<script setup>
import {useRoute, useRouter} from 'vue-router';
import {onMounted, ref} from "vue";
import { Bar } from 'vue-chartjs'
import { useTheme } from 'vuetify';
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  BarElement,
  CategoryScale,
  LinearScale
} from 'chart.js'


const router = useRouter();
const route = useRoute();
const theme = useTheme();
const userName = route.params.user;
const winCount = ref("");
const gameCount = ref("");
const looseCount = ref("");
const tryCountArr = ref([""]);
const avgTry = ref("");
const activeTheme = ref(theme.global.name.value);
ChartJS.register(Title, Tooltip, BarElement, CategoryScale, LinearScale);
const barChartData = ref({});
const barChartOptions = ref({});
const BarChart = Bar;


onMounted(() => {
  selectWinCount().then(data => winCount.value = data);
  selectGameCount().then(data => gameCount.value = data);
  selectLooseCount().then(data => looseCount.value = data);
  selectTryCount().then(data => tryCountArr.value = data.split(","));
  selectAVGTry().then(data => avgTry.value = data);
});

function getStat(id){
  if(id === "11") {
    let tmp = ['', '', '', '', '', ''];
    let j = 0;
    for (let i = 1; i <= 6; i++) {
      if (tryCountArr.value[j] === i.toString()) {
        tmp[i - 1] = tryCountArr.value[j + 1];
        j += 2;
      } else {
        tmp[i - 1] = 0;
      }
    }
    barChartData.value = {
      labels: ['1', '2', '3', '4', '5', '6'],
      datasets: [
        {
          label: 'Anzahl',
          data: [tmp[0], tmp[1], tmp[2], tmp[3], tmp[4], tmp[5]],
          backgroundColor: 'rgba(54, 162, 235, 0.5)',
          borderColor: 'rgba(54, 162, 235, 1)',
          borderWidth: 1
        }
      ]};
    barChartOptions.value = {
      indexAxis: 'y',
      responsive: true,
      scales: {
        x: {
          display: false // ← Versteckt die gesamte X-Achse
        },
        y: {
          grid: {
            display: false // ← Gitterlinien auf der Y-Achse ausblenden
          }
        }
      }
    };
  }
  switch(id){
    case "11": return "∑: "+ gameCount.value; break;
    case "12": return "✓: "+ winCount.value; break;
    case "13": return "✕: "+ looseCount.value; break;
    case "14": return "Ø: "+ avgTry.value; break;
  }
}
function getToolTip(id){
   switch(id){
    case "11": return "Gespielte Spiele"; break;
    case "12": return "Gewonnene Spiele"; break;
    case "13": return "Verlorene Spiele"; break;
    case "14": return "Durchschnittliche Versuche"; break;
  }
}

function routeBacktoMain() {
  router.push({
    name: 'home',
    params:{user:userName}
  });
}
function navigateLoginPage() {
  router.push({name: 'login'});
}

async function selectWinCount() {
  return fetch(`http://localhost:8080/api/selectWinCount?name=${encodeURIComponent(userName)}`)
      .then(response => response.text())
      .then(data => {
        return data;
      })
      .catch(error => {
        console.error("Fehler beim Abrufen:", error);
      });
}
async function selectGameCount() {
  return fetch(`http://localhost:8080/api/selectGameCount?name=${encodeURIComponent(userName)}`)
      .then(response => response.text())
      .then(data => {
        return data;
      })
      .catch(error => {
        console.error("Fehler beim Abrufen:", error);
      });
}
async function selectLooseCount() {
  return fetch(`http://localhost:8080/api/selectLooseCount?name=${encodeURIComponent(userName)}`)
      .then(response => response.text())
      .then(data => {
        return data;
      })
      .catch(error => {
        console.error("Fehler beim Abrufen:", error);
      });
}
async function selectTryCount() {
  return fetch(`http://localhost:8080/api/selectTryCount?name=${encodeURIComponent(userName)}`)
      .then(response => response.text())
      .then(data => {
        return data;
      })
      .catch(error => {
        console.error("Fehler beim Abrufen:", error);
      });
}
async function selectAVGTry() {
  return fetch(`http://localhost:8080/api/selectAVGTry?name=${encodeURIComponent(userName)}`)
      .then(response => response.text())
      .then(data => {
        return data;
      })
      .catch(error => {
        console.error("Fehler beim Abrufen:", error);
      });
}
</script>