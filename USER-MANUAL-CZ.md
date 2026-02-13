# PLZoREG Uživatelská Příručka

## Úvod

PLZoREG je automatický regulátor výkonu, který se instaluje mezi fotovoltaický (FV) střídač a elektrický ohřívač vody. Zabraňuje vypínání střídače nebo vzniku hluku během proměnlivého počasí tím, že plynule řídí výkon dodávaný do ohřívače.

![PLZoREG](cover.png)

## Instalace

**VAROVÁNÍ**: Instalace zahrnuje práci s vysokým napětím. Instalaci musí provést kvalifikovaný elektrikář v souladu s místními elektrotechnickými předpisy.

1. Namontujte zařízení na vhodné místo mezi střídač a ohřívač
2. Připojte vstup k AC výstupu vašeho FV střídače
3. Připojte výstup k vašemu ohřívači vody
4. Zajistěte správné uzemnění a ochranu obvodu
5. Zapněte systém

![alt text](drawing.png)

## Displej a Ovládání

Zařízení má **3-místný displej** a **dvě tlačítka**:
- **Tlačítko NAHORU** (↑)
- **Tlačítko DOLŮ** (↓)

Displej zobrazuje různé informace v závislosti na vybrané stránce.

## Provozní Stránky

Stisknutím tlačítka NAHORU nebo DOLŮ procházejte jednotlivé informační stránky:

### 1. Měřené Napětí
**Displej**: napětí bez desetinné tečky (např. `220`)

Zobrazuje aktuální napětí z vašeho FV střídače. Tato hodnota se mění podle oblačnosti.

- **Normální provoz**: 180-240V (závisí na vašem systému)
- **Akce**: Není vyžadována - pouze informativní

### 2. Cílové Napětí
**Displej**: napětí následované desetinnou tečkou (např. `220.`)

Toto je **cílové napětí**, které chcete udržovat. Zařízení sníží výkon do ohřívače, pokud napětí střídače klesne pod tuto hodnotu.

**Nastavení**:
1. Podržte OBĚ tlačítka po dobu 1 sekundy - desetinná tečka začne blikat
2. Stiskněte NAHORU pro zvýšení, DOLŮ pro snížení (rozsah: 100-240V)
3. Podržte OBĚ tlačítka po dobu 1 sekundy pro uložení a ukončení

**Doporučené nastavení**: 200-210V pro většinu systémů
- **Vyšší hodnoty** (220-230V): Konzervativnější, méně hluku, ale menší využití energie
- **Nižší hodnoty** (180-200V): Agresivnější, využívá více solární energie, ale může občas způsobit nestabilitu střídače

### 3. Pracovní Cyklus
**Displej**: Pracovní cyklus jako desetinné číslo (např. `0.75` = 75%)

Zobrazuje aktuální úroveň výkonu dodávaného do ohřívače:
- **0.00**: Žádný výkon (ohřívač vypnutý)
- **1.00**: Plný výkon (ohřívač na maximum)

Zařízení automaticky upravuje tuto hodnotu pro udržení cílového napětí.

**Ruční režim** (pokročilé):
1. Podržte OBĚ tlačítka po dobu 1 sekundy pro vstup do ručního režimu
2. Stiskněte NAHORU/DOLŮ pro úpravu po 5% krocích
3. Podržte OBĚ tlačítka pro ukončení ručního režimu (automatická regulace pokračuje)

### 4. Teplota chladiče triaku
**Displej**: `°` teplota následovaná symbolem ° (např. `50°`)

Zobrazuje teplotu spínacího prvku (informační hodnota):
- **Normální**: 20-80°C během provozu
- **Varování**: Nad 90°C spustí bezpečnostní vypnutí

### 5. Teplota MCU
**Displej**: `.°` teplota následovaná desetinnou tečkou a symbolem ° (např. `50.°`)

Zobrazuje vnitřní teplotu zařízení(informační hodnota):
- **Normální**: 20-60°C
- **Pouze informativní** - není vyžadována akce

## Chybové Kódy

Pokud dojde k chybě, displej zobrazí `E` následované číslem:

### E01 - Chybí Synchronizace
**Význam**: Zařízení nemůže detekovat synchronizační signál AC napájení

**Možné příčiny**:
- Chybí vstupní napájení
- Problém s kabeláží
- Hardwarová závada

**Řešení**: Zkontrolujte, že střídač je zapnutý a vydává AC napětí. Pokud problém přetrvává, kontaktujte podporu.

### E02 - Přehřátí
**Význam**: Teplota triaku překročila 90°C

**Možné příčiny**:
- Porucha termostatu ohřívače
- Nedostatečný průtok/cirkulace vody
- Provoz nasucho (žádná voda v ohřívači)

**Řešení**:
1. Okamžitě odpojte napájení
2. Zkontrolujte hladinu vody a cirkulaci v ohřívači
3. Ověřte funkci termostatu ohřívače
4. Před restartováním nechte systém vychladnout

**Bezpečnostní upozornění**: Když je aktivní jakákoli chyba, výstup výkonu je automaticky vypnutý.

## Normální Provoz

1. **Slunečné podmínky**:
   - Displej zobrazuje V.sense blízko nebo nad V.target
   - Pracovní cyklus se blíží 100%
   - Maximální výkon dodávaný do ohřívače

2. **Částečně zataženo**:
   - Displej zobrazuje kolísající V.sense
   - Pracovní cyklus se automaticky upravuje (20-90%)
   - Zařízení reguluje výkon plynule

3. **Velmi zataženo/málo světla**:
   - Displej zobrazuje V.sense na nebo pod V.target
   - Pracovní cyklus blízko 0%
   - Minimální výkon do ohřívače

## Tipy pro Optimální Provoz

- **Nastavte cílové napětí konzervativně** při prvním použití - začněte na 210V a snižujte, pokud je systém stabilní
- **Sledujte několik dní** za různých povětrnostních podmínek před finálním nastavením
- **Pravidelně kontrolujte teplotní stránku**, abyste se ujistili, že ohřívač funguje normálně
- **Nenastavujte cíl příliš nízko** - tím ztrácíte smysl a můžete způsobit hluk, kterému se snažíte vyhnout

## Údržba

- Není vyžadována pravidelná údržba
- Udržujte ventilační otvory čisté
- Pravidelně kontrolujte, že displej funguje
- Ověřte, že teplotní údaje jsou rozumné

## Specifikace

- **Vstupní napětí**: 100-240V AC (z FV střídače)
- **Výstupní řízení**: 0-100% fázové PWM
- **Rozsah cílového napětí**: 100-240V (uživatelsky nastavitelné)
- **Sledování teploty**: Až do 90°C (bezpečnostní limit)
- **Reakční čas**: 500ms aktualizační cyklus
- **Uchování nastavení**: Nevolatilní (přežije ztrátu napájení)

## Řešení Problémů

| Příznak | Možná Příčina | Řešení |
|---------|---------------|--------|
| Displej nesvítí | Chybí napájení | Zkontrolujte vstupní napájení |
| Trvalá chyba E01 | Chybí AC vstup | Ověřte, že střídač je aktivní |
| Častá chyba E02 | Problém s ohřívačem | Zkontrolujte termostat/hladinu vody ohřívače |
| Střídač stále hlučí | Cílové napětí příliš nízké | Zvyšte nastavení cílového napětí |
| Nízké využití energie | Cílové napětí příliš vysoké | Snižte nastavení cílového napětí |
| Displej zobrazuje špatné hodnoty | Vyžaduje kalibraci | Kontaktujte podporu |

## Technická Podpora

[Doplňte kontaktní informace]

## Bezpečnostní Varování

- ⚠️ **VYSOKÉ NAPĚTÍ** - Neotvírejte kryt zařízení
- ⚠️ **Instalace pouze kvalifikovaným personálem**
- ⚠️ **Nepřemosťujte chybové stavy**
- ⚠️ **Zajistěte správné uzemnění**
