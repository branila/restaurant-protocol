# SBURP - Simple But Useful Restaurant Protocol

## Descrizione Generale

SBURP (Simple But Useful Restaurant Protocol) è un protocollo di formattazione testuale progettato per gestire comande e ordini in ambito ristorativo. Questo sistema permette di tradurre in modo efficiente richieste culinarie dalla forma testuale a formati strutturati, facilitando la comunicazione tra la sala e la cucina.

## Perché SBURP?

Perché complicare ciò che può essere semplice? SBURP trasforma le ordinazioni caotiche in strutture ordinate.

## Punti di Forza

- **Semplicità**: Sintassi intuitiva che richiede minima formazione per il personale
- **Flessibilità**: Supporta modifiche ai piatti senza compromettere la struttura dell'ordine
- **Compatibilità**: Converte facilmente in vari formati (JSON, XML, YAML) per integrarsi con diversi sistemi
- **Validazione Integrata**: Verifica automaticamente gli errori prima di inviare l'ordine in cucina
- **Leggibilità**: Output formattato comprensibile sia dalle macchine che dagli umani
- **Zero Dipendenze**: Non richiede software specializzato per essere implementato

## Sintassi

La sintassi SBURP è strutturata gerarchicamente, permettendo di definire ordini, comande e singoli piatti con le loro modifiche.

### Struttura Base

```
ORDINE <ID_ORDINE> <DATA>

COMANDA <ID_COMANDA>
<PORTATA> "<NOME_PIATTO>" [+"<INGREDIENTE>"] [-"<INGREDIENTE>"]

COMANDA <ID_COMANDA>
<PORTATA> "<NOME_PIATTO>" [+"<INGREDIENTE>"] [-"<INGREDIENTE>"]
...
```

### Dettagli Sintattici

- **ORDINE**: Definisce il numero del tavolo e la data dell'ordine
- **COMANDA**: Raggruppa più piatti ordinati contemporaneamente, identificati da un numero progressivo
- **Tipo Piatto**: Può essere `PRIMO`, `SECONDO` o `CONTORNO`
- **Nome Piatto**: Sempre racchiuso tra virgolette doppie (`"..."`)
- **Modifiche**: Indicate con `+` (aggiunta) o `-` (rimozione) seguito dal nome dell'ingrediente tra virgolette

## Esempi

### Input Valido

```
ORDINE 1 13/11/2025

COMANDA 0
PRIMO "pasta al pomodoro" +"formaggio" -"basilico"
CONTORNO "insalata" -"olio" +"aceto balsamico"

COMANDA 1
PRIMO "pasta al pomodoro" -"pomodoro"
SECONDO "Bistecca" +"Salsa barbecue"
CONTORNO "insalata" -"olio"
```

### Output Formattato

```
Ordine per il tavolo 1, data 13/11/2025
  Comanda 0:
    Primo: pasta al pomodoro [{+ formaggio} {- basilico}]
    Contorno: insalata [{- olio} {+ aceto balsamico}]
  Comanda 1:
    Primo: pasta al pomodoro [{- pomodoro}]
    Secondo: Bistecca [{+ Salsa barbecue}]
    Contorno: insalata [{- olio}]
```

## Formati di Conversione Supportati

SBURP supporta la conversione in vari formati per facilitare l'integrazione con sistemi diversi.

### JSON

```json
{
  "Tavolo": 1,
  "Data": "13/11/2025",
  "Comande": [
    {
      "Numero": 0,
      "Primo": {
        "Nome": "pasta al pomodoro",
        "Modifiche": [
          {
            "Tipo": "+",
            "Voce": "formaggio"
          },
          {
            "Tipo": "-",
            "Voce": "basilico"
          }
        ]
      },
      "Secondo": null,
      "Contorno": {
        "Nome": "insalata",
        "Modifiche": [
          {
            "Tipo": "-",
            "Voce": "olio"
          },
          {
            "Tipo": "+",
            "Voce": "aceto balsamico"
          }
        ]
      }
    },
    {
      "Numero": 1,
      "Primo": {
        "Nome": "pasta al pomodoro",
        "Modifiche": [
          {
            "Tipo": "-",
            "Voce": "pomodoro"
          }
        ]
      },
      "Secondo": {
        "Nome": "Bistecca",
        "Modifiche": [
          {
            "Tipo": "+",
            "Voce": "Salsa barbecue"
          }
        ]
      },
      "Contorno": {
        "Nome": "insalata",
        "Modifiche": [
          {
            "Tipo": "-",
            "Voce": "olio"
          }
        ]
      }
    }
  ]
}
```

### XML

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Ordine>
  <Tavolo>1</Tavolo>
  <Data>13/11/2025</Data>
  <Comande>
    <Comanda Numero="0">
      <Primo>
        <Nome>pasta al pomodoro</Nome>
        <Modifiche>
          <Modifica>
            <Tipo>+</Tipo>
            <Voce>formaggio</Voce>
          </Modifica>
          <Modifica>
            <Tipo>-</Tipo>
            <Voce>basilico</Voce>
          </Modifica>
        </Modifiche>
      </Primo>
      <Contorno>
        <Nome>insalata</Nome>
        <Modifiche>
          <Modifica>
            <Tipo>-</Tipo>
            <Voce>olio</Voce>
          </Modifica>
          <Modifica>
            <Tipo>+</Tipo>
            <Voce>aceto balsamico</Voce>
          </Modifica>
        </Modifiche>
      </Contorno>
    </Comanda>
    <Comanda Numero="1">
      <Primo>
        <Nome>pasta al pomodoro</Nome>
        <Modifiche>
          <Modifica>
            <Tipo>-</Tipo>
            <Voce>pomodoro</Voce>
          </Modifica>
        </Modifiche>
      </Primo>
      <Secondo>
        <Nome>Bistecca</Nome>
        <Modifiche>
          <Modifica>
            <Tipo>+</Tipo>
            <Voce>Salsa barbecue</Voce>
          </Modifica>
        </Modifiche>
      </Secondo>
      <Contorno>
        <Nome>insalata</Nome>
        <Modifiche>
          <Modifica>
            <Tipo>-</Tipo>
            <Voce>olio</Voce>
          </Modifica>
        </Modifiche>
      </Contorno>
    </Comanda>
  </Comande>
</Ordine>
```

### YAML

```yaml
Tavolo: 1
Data: 13/11/2025
Comande:
  - Numero: 0
    Primo:
      Nome: pasta al pomodoro
      Modifiche:
        - Tipo: "+"
          Voce: formaggio
        - Tipo: "-"
          Voce: basilico
    Secondo: null
    Contorno:
      Nome: insalata
      Modifiche:
        - Tipo: "-"
          Voce: olio
        - Tipo: "+"
          Voce: aceto balsamico
  - Numero: 1
    Primo:
      Nome: pasta al pomodoro
      Modifiche:
        - Tipo: "-"
          Voce: pomodoro
    Secondo:
      Nome: Bistecca
      Modifiche:
        - Tipo: "+"
          Voce: Salsa barbecue
    Contorno:
      Nome: insalata
      Modifiche:
        - Tipo: "-"
          Voce: olio
```

## Gestione Errori

SBURP include un sistema di validazione che verifica la correttezza delle richieste. In caso di incongruenze, viene generato un messaggio di errore che specifica:

- Codice errore
- Descrizione del problema
- Possibile soluzione
- Suggerimento per l'utente

### Esempio di Errore

Input con piatto non esistente:
```
ORDINE 1 13/11/2025
COMANDA 0
PRIMO "pasta alle carote" +"formaggio" -"basilico"
...
```

Output di errore:
```
=== ERRORE [1002] ===
Messaggio: Il piatto 'pasta alle carote' non esiste nel menu
Soluzione: Consultare il menu aggiornato per verificare i piatti disponibili
Suggerimento: Controllare il menu per piatti alternativi disponibili.
```

## Codici di Errore Comuni

| Codice | Descrizione |
|--------|-------------|
| 1001   | Formato SBURP non valido |
| 1002   | Piatto non esistente nel menu |
| 1003   | Modifica non valida per il piatto specificato |
| 1004   | Data non valida |
| 1005   | Numero tavolo non valido |

## Note Implementative

- Le modifiche vengono processate nell'ordine in cui appaiono nel testo
- Una modifica che rimuove un ingrediente non presente nel piatto base genererà un errore
- SBURP è case-sensitive: prestare attenzione a maiuscole e minuscole
- L'output formattato mantiene l'ordine delle modifiche come specificato nell'input

## Integrazione

SBURP può essere facilmente integrato con:
- Sistemi POS (Point of Sale)
- Tablet e dispositivi mobili per la presa di ordinazioni
- Sistemi di gestione cucina (KDS)
- Software di fatturazione e inventario
- App di delivery

## Conclusione

SBURP è come un buon sous-chef: fa il lavoro pesante senza lamentarsi, permettendovi di concentrarvi sull'essenziale. 
Implementate SBURP nel vostro sistema di gestione ristorativa e osservate come le comande fluiscono dalla sala alla cucina con l'eleganza di un sommelier che versa un Barolo d'annata.

*"Se non è SBURP, è probabilmente troppo complicato."*
