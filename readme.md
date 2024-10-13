# Go-tutoriaali Projekti

Tämä projekti on Go-pohjainen API, joka demonstroi CRUD-toimintoja MongoDB:n avulla sekä edistyneitä ominaisuuksia, kuten sivutusta, lajittelua ja suodatusta. API on dokumentoitu Swaggerilla ja se on otettu käyttöön Dockerin avulla Azureen.

## Tutoriaali Osuus
Tein ensin kaksi ensimmäistä Go-tutoriaalia ilman, että mietin vaihetta 2. Kolmannen tutoriaalin aloitin alusta asti vaihe 2 mielessä käyttäen MongoDB:tä, kuten tässä kurssissa tehtiin.

## API:n Toteutus
API sisältää seuraavat päätepisteet `Category`, `Animals` ja `Species`:
- **GET**: Hakee tieto kannasta joko kaikki tai yksittäisen
- **POST**: Lisää uuden tiedon tietokantaan
- **PUT**: Päivittää tietokannasta jo löytyvää tietoa
- **DELETE**: Poistaa tietokannasta

`GET /animals` päätepisteessä on käytetty `populate`-toimintoa, jotta saadaan `category` ja `species` tiedot mukaan.

## Edistyneet Ominaisuudet
- **Sivutus**: Toteutettu kaikissa `Category`, `Animals` ja `Species` päätepisteissä. Tällä määrätää sivun määrä tietoihin
- **Lajittelu**: Lisätty samoihin päätepisteisiin kuin sivutus. Tällä voidaan määrittää m issä järjestyksessä tiedot tulevat
- **Suodatus**: Toteutettu `Species` päätepisteessä, jossa voi hakea tiettyä lajin nimellä.

## Dokumentointi
API-dokumentaatio on luotu Swaggerin avulla. Tämä sisältää tietoa API-reiteistä

## Käyttöönotto
Projekti on otettu käyttöön Azureen Docker-konttin avulla.

## Suoritetut Vaiheet
1. Tein Go-tutoriaalit.
2. Käytän tehtävissä tehtyä MongoDB:tä ja kaikki perus CRUD-toiminnot ovat käytössä.
3. Käytän Filtering, Sorting ja Pagination -toimintoja. Esimerkistä pääsee kokeilemaan suodatusta.
4. Käytän Swaggeria dokumentoinnissa.
5. Laitettu Dockeriin ja pyöritetään Azure:ssa.