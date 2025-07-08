
<template>
    <v-app-bar app>
        <v-toolbar-title> Wordle IBT </v-toolbar-title>

        <v-btn icon class="mr-4" density="compact" @click="{activeTheme == 'light' ? activeTheme = 'dark' : activeTheme ='light'; theme.global.name.value = activeTheme}">
            <v-icon v-if="activeTheme == 'light'"> mdi-white-balance-sunny </v-icon>
            <v-icon v-if="activeTheme == 'dark'"> mdi-weather-night </v-icon>
        </v-btn>
        <p class="mr-4" disabled>|</p>
        <v-btn class="mr-4" density="compact" icon="mdi-poll" @click="navigateStatsPage"> </v-btn>
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
            <div class="justify-center autofit-width">
                <v-alert
                    v-if="alertVisible"
                    type="warning"
                    color="orange"
                    class="mt-4 mb-4"
                    variant="tonal"
                    >
                    {{alertText}}
                </v-alert>
            </div>
        </v-container>

        <v-container class="d-flex justify-center pa-0">
        <v-sheet class="border pa-8 ma-0" width="900" elevation="20">
            <div class="d-flex justify-center">
                <v-btn 
                    id="newGameBtn"
                    v-if="levelCompleted"
                    class="mb-8"
                    color="primary" 
                    @click="resetGame"
                >  
                    New Game
                </v-btn>
            </div>
            <v-container class="autofit-width pa-0">
                <v-row v-for="m in 6" :key="m" class="autofit-width autofit-margin border rounded pa-1">
                    <div v-for="n in 5" :key="n" class="ma-1" color="theme_color_1">
                        <v-text-field
                            density="compact"
                            class="wordle-cell"
                            variant="outlined" 
                            maxlength="1" 
                            autocomplete="off"
                            hide-details
                            v-model="grid[m-1][n-1]"
                            :id="`${m}${n}`" 
                            :readonly="currentRow != m || levelCompleted || denieInput"
                            :disabled="currentRow < m && currentRow != 0"
                            @keydown="handleKeydown" 
                            @update:model-value="handleKeyInput(`${m}${n}`)"
                        ></v-text-field >
                    </div>
                </v-row>

                <div class="autofit-width autofit-margin mt-8">
                    <v-text-field 
                        v-if="levelCompleted && currentRow == 0"
                        class="wordle-cell color-green"
                        variant="outlined" 
                        width="200"
                        hide-details
                        readonly
                        :model-value="upperSWord"
                    ></v-text-field>
                </div>
            </v-container>

            <v-container class="autofit-width border rounded mt-8 elevation-10">
                <v-row
                v-for="(row, rowIndex) in keyboardRows"
                :key="rowIndex"
                class="justify-center"
                >
                <v-btn
                    v-for="(key) in row"
                    :id="key.toLowerCase()"
                    class="ma-1 pl-6 pr-6 pt-5 pb-8"
                    :class="getKeyColorClass(key)"
                    variant="tonal" 
                    @mousedown.prevent="pressLetterBtn(key)"
                >
                    {{ key }}
                </v-btn>
                </v-row>
            </v-container>

        <v-dialog id="endDialog" v-model="showEndDialog" width="500" persistent>
          <v-card>
              <v-card-title v-if="currentRow != 0">
                  Level erfolgreich Abgeschlossen!
              </v-card-title>
              <v-card-title v-if="currentRow == 0">
                  Level nicht geschafft.
              </v-card-title>
              <v-card-text>
                <template v-if="currentRow == 1">
                  <p> {{ `Du hast einen Versuch gebraucht` }} </p>
                </template>
                <template v-if="currentRow > 1">
                  <p> {{ `Du hast ${currentRow} Versuche gebraucht` }} </p>
                </template>
                <template v-if="currentRow == 0">
                  <p> {{ `Das Wort war: ${solutionWord.toUpperCase()}` }} </p>
                </template>
                  <v-spacer></v-spacer>
                  <p v-if="currentRow != 0"> {{ `Neuer Durchschnitt: ${avgTry} ` }}</p>
              </v-card-text>
              <v-card-actions>
                  <v-btn @click="resetGameByDialog">
                    New Game
                  </v-btn>
                  <v-btn @click="endDialogClose">
                    Close
                  </v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
        </v-sheet>
      </v-container>
    </v-main>
</template>


<script setup>
  
    import { onMounted, ref, watch, nextTick } from "vue";
    import { useRouter, useRoute } from 'vue-router';
    import { useTheme } from 'vuetify';
    import { wordList } from "../../words/random-word-api-cleaned";

    const router = useRouter();
    const route = useRoute();
    const userName = route.params.user; // importiert aus route url
    const theme = useTheme();
    const activeTheme = ref(theme.global.name.value);

    const urlBase = 'https://de.wiktionary.org/w/api.php';
    //const urlBase2 = 'https://random-word-api.herokuapp.com' => lokal als Datei anstatt API call

    const levelCompleted = ref(false);
    const showEndDialog = ref(false);
    const denieInput = ref(false);
    const currentRow = ref(1);
    const avgTry = ref("");
    const alertText = ref("");
    const alertVisible = ref(false);

    const upperSWord = ref('');
    let solutionWord = '';

    const grid = ref([
        ['','','','',''],
        ['','','','',''],
        ['','','','',''],
        ['','','','',''],
        ['','','','',''],
        ['','','','','']
    ]);

     const keyboardRows= [
        ["Q", "W", "E", "R", "T", "Z", "U", "I", "O", "P"],
        ["A", "S", "D", "F", "G", "H", "J", "K", "L"],
        ["EINGABE", "Y", "X", "C", "V", "B", "N", "M", "LÖSCHEN"]
    ];
    const keyColors = ref({});

    onMounted(() => {
        generateRandomWord();
        setElementFocusById("11");
        selectAVGTry().then(data => avgTry.value = data);
        window.addEventListener('beforeunload', handleBeforeUnload);
    });
    
    watch(showEndDialog, (val) => {
        if (!val) {
            nextTick(() => {
            setElementFocusById("11");
            })
        }
    });

    const handleBeforeUnload = (event) => {
        if(!levelCompleted.value && currentRow.value != 1){
            event.preventDefault();
            event.returnValue = '';
        }
    };

    function generateRandomWord() {

        levelCompleted.value = false;

        const listLength = wordList.length;
        const randNum = Math.random()*(listLength-1);

        solutionWord = wordList.at(randNum).toLowerCase();
        upperSWord.value = solutionWord.split('').join(' ').toUpperCase();
        console.log("---");
        console.log("SolutionWord: " + solutionWord);
        console.log("---");
    }
    function handleKeyInput(textfieldId){
        
        const id = parseInt(textfieldId);
        const field = document.getElementById(id);

        // wenn gültiges zeichen eingegeben wurde, ein feld weiter
        if(field.value != ''){
            const nextfield = document.getElementById(id + 1);
            nextfield?.focus();
        }

        grid.value[textfieldId[0]-1][textfieldId[1]-1] = grid.value[textfieldId[0]-1][textfieldId[1]-1].toUpperCase();
    }
    async function handleKeydown(event) {

        const fieldId = event.target.id;
        const intId = parseInt(event.target.id);
        const field = document.getElementById(fieldId);

        if(levelCompleted.value|| fieldId[0] != currentRow.value || denieInput.value) return;

        if(event.key === 'Enter'){
            denieInput.value = true;
            handlePressEnter();
            return;
        }

        // Pfeiltasten-navigation ermöglichen
        // Und taste unterdrücken damit Cursor-position nicht vor dem Buchstaben
        if(event.key === 'ArrowLeft'){
            let nextfield = document.getElementById(intId - 1);
            nextfield?.focus();
            event.preventDefault();
            return;
        }
        if(event.key === 'ArrowRight'){
            let nextfield = document.getElementById(intId + 1);
            nextfield?.focus();
            event.preventDefault();
            return;
        }

        // wenn Backspace gedrückt wird und das text-field leer ist, ein feld zurück
        // wenn Backspace gedrückt wird und das text-field nicht leerist, im feld bleiben
        if (event.key === 'Backspace') {
            event.preventDefault();
            handlePressDelete(field, fieldId);
            return;
        }

        // ungültige tasten unterdrücken
        else if (!/^[a-z]$/i.test(event.key)) {
            event.preventDefault();
            return;
        }
        else{
            // Damit man Buchstaben überschreiben kann
            grid.value[fieldId[0]-1][fieldId[1]-1] = '';
        }
    }
    async function handleWordInput(){

        const guessWord = getGuessedWord();
        const capitalGuessWord = guessWord.charAt(0).toUpperCase() + guessWord.substring(1);

        const exists1 = await checkWordExist(guessWord).valueOf();
        const exists2 = await checkWordExist(capitalGuessWord).valueOf();

        console.log("exists: "+ guessWord + " - " + exists1);
        console.log("exists: "+ capitalGuessWord + " - " + exists2);

        if(!exists1 && !exists2){
            alertText.value = "Das Wort ist ungültig";
            showAlert();
            console.log("Worteingabe angelehnt: Das Wort ist ungültig");
            denieInput.value = false;
            return;
        }

        compareWords(guessWord);  

        // Spiel Ende
        if(guessWord == solutionWord || currentRow.value == 6){

            if(guessWord != solutionWord) currentRow.value = 0;

            levelCompleted.value = true;
            db_insert(userName, currentRow.value.toString(), solutionWord);
            sleep(100);
            reloadAVG()
            showEndDialog.value = true;
            denieInput.value = false;
            return;
        }

        // fokus auf nächste reihe
        currentRow.value++;
        let nextfield = document.getElementById(`${currentRow.value}${1}`);
        nextfield.disabled = false;
        denieInput.value = false;
        nextfield?.focus();
    }
    
    /*  CompareWords Methode:
    *   - Das GuessWord mit dem Solutionword vergleichen und die UI-Elemente einfärben.
    *   
    *   Vorgehen:
    *   - arrayGWord iterieren und auf 'grüne' Buchstaben prüfen
    *   - Wenn gefunden: Chars in beiden Arrays mit '#' ersetzen und UI grün färben
    *   
    *   - arrayGWord iterieren und auf 'orangene' Buchstaben prüfen
    *   - Dabei Chars in arrayGWord auslassen, die == '#', weil die sind schon grün
    *   - Für jeden anderen Char durch arraySWord iterieren
    *   - Wenn gefunden: Chars in beiden Arrays mit '#' ersetzen und UI orange färben
    *   - Wenn nicht: Chars in Arrays beibehalten  und UI grau färben
    */
    function compareWords(guessdWord) {

        let arraySWord = solutionWord.split('');
        let arrayGWord = guessdWord.split('');

        for(let i=0; i<5; i++){

            // if Buchstaben an richtiger Stelle
            if(arrayGWord[i] == arraySWord[i]){
                const id = `${currentRow.value}${i+1}`;
                setElementColorById(id, 'green');
                updateKeyColor(arrayGWord[i], 'green');
                arrayGWord[i] = '#';
                arraySWord[i] = '#';
            }
        }
        for(let i=0; i<5; i++){ // arrayGWord iterieren

            if(arrayGWord[i] != '#'){ // damit grüne Buchstaben nicht geprüft werden

                for(let j=0; j<5; j++){ // arraySWord iterieren

                    // if Buchstabe an falscher Stelle
                    if(arrayGWord[i] == arraySWord[j]){
                        const id = `${currentRow.value}${i+1}`;
                        setElementColorById(id, 'orange');
                        updateKeyColor(arrayGWord[i], 'orange');
                        arrayGWord[i] = '#';
                        arraySWord[j] = '#';
                        break;
                    }
                }
                // if Buchstabe kommt nicht vor
                if(arrayGWord[i] != '#'){
                    const id = `${currentRow.value}${i+1}`;
                    setElementColorById(id, 'grey');
                    updateKeyColor(arrayGWord[i], 'grey');
                }
            }
        }
    }

    function getGuessedWord() {
        // buchstaben aus textfeldern in string packen

        let guessWord = grid.value[currentRow.value-1].join('');
        guessWord = guessWord.toLowerCase();
        console.log("GuessWord: " + guessWord)
        
        return guessWord.toLowerCase();
    }

    function setElementFocusById(id) {
        const field = document.getElementById(id);
        field?.focus();
    }

    function setElementColorById(id, color) {

        const field = document.getElementById(id);

        if(color == 'green' && !field.classList.contains('color-grey')){
            field.classList.remove('color-orange');
            field.classList.add(`color-${color}`);
        }
        else if(color == 'orange' && !field.classList.contains('color-green') && !field.classList.contains('color-grey')){
            field.classList.add(`color-${color}`);
        }
        else if(color == 'grey' && !field.classList.contains('color-orange') && !field.classList.contains('color-green')){
            field.classList.add(`color-${color}`);
        }
    }
    async function checkWordExist(word){
        try{
            const response = await fetch(`${urlBase}?action=query&titles=${encodeURIComponent(word)}&format=json&origin=*`);
            const data = await response.json();
            const pages = data.query.pages;
            const exists = Object.keys(pages)[0] !== "-1";
            return exists;
        } catch (error){
            console.log("Fehler beim Abrufen der Word-API:", error);
            return true;
        }
    }

    function pressLetterBtn(btnId) {

        const field = document.activeElement;
        const fieldId = field.id;
        const intId = parseInt(fieldId);

        if(levelCompleted.value|| fieldId[0] != currentRow.value) return;

        if(btnId == 'EINGABE'){
            handlePressEnter();
            return;
        }
        if(btnId == 'LÖSCHEN'){
            handlePressDelete(field, fieldId);
            return;
        }

        grid.value[fieldId[0]-1][fieldId[1]-1] = btnId;
        field.value = btnId;

         if(field.value != ''){
            const nextfield = document.getElementById(intId + 1);
            nextfield?.focus();
        }
    }

    function getKeyColorClass(key) {
        const lowerKey = key.toLowerCase(); 
        const color = keyColors.value[lowerKey];
        return color ? `color-${color}` : '';
    }

    function updateKeyColor(key, color) {
        const priority = { grey: 1, orange: 2, green: 3 };
        const current = keyColors.value[key];
        if (!current || color == '' || priority[color] > priority[current]) {
            keyColors.value[key] = color;
        }
    }

    function handlePressEnter(){
        if( grid.value[currentRow.value-1][0] != '' && 
            grid.value[currentRow.value-1][1] != '' && 
            grid.value[currentRow.value-1][2] != '' && 
            grid.value[currentRow.value-1][3] != '' && 
            grid.value[currentRow.value-1][4] != ''
        ){
            handleWordInput();
        }
        else{
            alertText.value = "Nicht alle Textfelder sind befüllt";
            showAlert();
            denieInput.value = false;
        }
    }
    
    function handlePressDelete(field, id){

        if(field.value == '') {

            let nextfield = document.getElementById(id - 1);
            nextfield?.focus();
            grid.value[id[0]-1][id[1]-2] = '';
            return;
        }
        else{
            grid.value[id[0]-1][id[1]-1] = '';
            return; 
        }
    }

    function resetGameByDialog(){
      showEndDialog.value = false;
      resetGame();
    }

    function endDialogClose() {
        showEndDialog.value = false;
        
    }

    function resetGame() {
        for(let i=1; i<=6; i++) {
            for(let j=1; j<=5; j++) {
                let fieldId = i.toString() + j;
                let field = document.getElementById(fieldId);
                grid.value[i-1][j-1] = '';
                field.classList.remove('color-orange', 'color-green', 'color-grey');
            }
        }
       
        keyboardRows.flat().forEach((key) =>{
            updateKeyColor(key.toLowerCase(), '')
        })
        currentRow.value = 1;
        denieInput.value = false;
        
        generateRandomWord();
        let nextfield = document.getElementById("11");
        nextfield.disabled = false;
        nextfield?.focus();
    }

    function reloadAVG(){
      selectAVGTry().then(data => avgTry.value = data);
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

    async function db_insert(name, versuche, wort) {

        try{
            await fetch("http://localhost:8080/api/insert", { 
                method: "POST", 
                body: JSON.stringify({ name: name, versuche: versuche, wort: wort})
            });
        }
        catch(error){
            console.log("Fehler beim Abrufen oder Verarbeiten der Daten:", error);
        }
    }

    function sleep(milliseconds) {
      var start = new Date().getTime();
      for (var i = 0; i < 1e7; i++) {
        if ((new Date().getTime() - start) > milliseconds){
        break;
        }
      }
    }

    function showAlert() {
        alertVisible.value = true

        // Nach 2 Sekunden automatisch ausblenden
        setTimeout(() => {
            alertVisible.value = false
        }, 2000)
    }
    function navigateStatsPage() {
        if(!levelCompleted.value && currentRow.value != 1){
            let tmp = confirm("Sie haben ungespeicherte Änderungen.\nWenn sie Fortfahren geht Ihr Fortschritt Verloren!");
            if(!tmp) return;
        }
        router.push({
        name: 'stats',
        params:{user:userName}
        });
    }
    function navigateLoginPage() {
        if(!levelCompleted.value && currentRow.value != 1){
            let tmp = confirm("Sie haben ungespeicherte Änderungen.\nWenn sie Fortfahren geht Ihr Fortschritt Verloren!");
            if(!tmp) return;
        }
        router.push({name: 'login'});
    }
</script>
