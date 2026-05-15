;; GENESIS.DSL — The Creation Program
;; A Hebrew-native executable specification
;; Run with: go run ./experiments/dsl   (auto-locates this file; -fast skips pauses)

;; ══════════════════════════════════════════════════════════════
;; HEADER — בְּרֵאשִׁית (In the beginning)
;; ══════════════════════════════════════════════════════════════

(program "genesis"
  :version "1.0.0"
  :author "Moses (traditional)"
  :encoding "UTF-8"
  :direction "RTL")

;; ══════════════════════════════════════════════════════════════
;; 1:1-2 — INITIALIZATION
;; ══════════════════════════════════════════════════════════════

(bara ELOHIM                          ;; God created (ex nihilo)
  (et HA-SHAMAYIM)                    ;; the heavens
  (ve-et HA-ARETZ))                   ;; and the earth

(define ERETZ
  :state   (TOHU VA-VOHU)             ;; formless and void
  :surface (CHOSHEKH AL-PNEI TEHOM)   ;; darkness on face of deep
  :hovering (RUACH ELOHIM AL-PNEI HA-MAYIM))  ;; Spirit over waters

;; ══════════════════════════════════════════════════════════════
;; 1:3-5 — DAY 1: LIGHT
;; ══════════════════════════════════════════════════════════════

(yomer ELOHIM                         ;; God said
  (yehi OR))                          ;; "Let there be light"
(va-yehi OR)                          ;; And there was light → OUTPUT

(va-yar ELOHIM et HA-OR ki TOV)       ;; God saw light: GOOD
(va-yavdel ELOHIM                     ;; God separated
  :between HA-OR
  :and     HA-CHOSHEKH)

(va-yikra ELOHIM                      ;; God called/named
  (la-or    YOM)                      ;; light → "Day"
  (la-choshekh LAILAH))               ;; darkness → "Night"

(va-yehi EREV va-yehi VOKER           ;; evening, morning
  :day 1)                             ;; Day One (אֶחָד)

;; ══════════════════════════════════════════════════════════════
;; 1:6-8 — DAY 2: SKY/FIRMAMENT
;; ══════════════════════════════════════════════════════════════

(yomer ELOHIM
  (yehi RAKIA be-tokh HA-MAYIM)       ;; firmament in waters
  (vihi MAVDIL                        ;; let it divide
    :between MAYIM
    :and     MAYIM))

(va-yaas ELOHIM et HA-RAKIA           ;; God made firmament
  (va-yavdel
    :between (MAYIM mi-tachat la-RAKIA)   ;; waters below
    :and     (MAYIM me-al    la-RAKIA)))  ;; waters above
(va-yehi KEN)                         ;; And it was so → TRUE

(va-yikra ELOHIM la-RAKIA SHAMAYIM)   ;; called firmament "Sky"

(va-yehi EREV va-yehi VOKER :day 2)

;; ══════════════════════════════════════════════════════════════
;; 1:9-13 — DAY 3: LAND, SEAS, PLANTS
;; ══════════════════════════════════════════════════════════════

(yomer ELOHIM
  (yikavu HA-MAYIM                    ;; let waters gather
    :to (MAKOM ECHAD))                ;; to one place
  (ve-teraeh HA-YABASHAH))            ;; let dry land appear
(va-yehi KEN)

(va-yikra ELOHIM
  (la-yabashah ERETZ)                 ;; dry land → "Earth"
  (u-le-mikveh ha-mayim YAMIM))       ;; gathered waters → "Seas"
(va-yar ELOHIM ki TOV)

;; COGNATE STRUCTURE: תַּדְשֵׁא...דֶּשֶׁא (vegetate vegetation)
(yomer ELOHIM
  (tadshe HA-ERETZ DESHE              ;; earth.vegetate(vegetation)
    (ESEV mazria ZERA)                ;; plant.seed(seed)
    (ETZ PRI oseh PRI                 ;; tree.fruit(fruit)
      :le-mino   SELF                 ;; after its kind
      :zaro-vo   SELF)))              ;; seed in itself
(va-yehi KEN)

(va-totze HA-ERETZ DESHE              ;; earth brought forth
  (ESEV mazria ZERA le-minehu)
  (ve-ETZ oseh PRI asher zaro-vo le-minehu))
(va-yar ELOHIM ki TOV)

(va-yehi EREV va-yehi VOKER :day 3)

;; ══════════════════════════════════════════════════════════════
;; 1:14-19 — DAY 4: SUN, MOON, STARS
;; ══════════════════════════════════════════════════════════════

(yomer ELOHIM
  (yehi MEOROT bi-RAKIA HA-SHAMAYIM   ;; lights in firmament
    :le-havdil (bein HA-YOM u-vein HA-LAILAH)
    :le-otot   (MOADIM YAMIM SHANIM)))  ;; signs, seasons, days, years
(va-yehi KEN)

(va-yaas ELOHIM
  (et HA-MAOR HA-GADOL :le-memshelet HA-YOM)    ;; greater light → day
  (et HA-MAOR HA-KATAN :le-memshelet HA-LAILAH) ;; lesser light → night
  (ve-et HA-KOKHAVIM))                           ;; and the stars
(va-yar ELOHIM ki TOV)

(va-yehi EREV va-yehi VOKER :day 4)

;; ══════════════════════════════════════════════════════════════
;; 1:20-23 — DAY 5: SEA CREATURES, BIRDS
;; ══════════════════════════════════════════════════════════════

(yomer ELOHIM
  (yishretzu HA-MAYIM SHERETZ NEFESH CHAYAH)  ;; waters swarm
  (ve-OF yeofeif al-HA-ARETZ))                 ;; birds fly
(va-yehi KEN)

(va-yivra ELOHIM                               ;; God CREATED (bara!)
  (et HA-TANINIM HA-GEDOLIM)                   ;; great sea creatures
  (ve-et kol NEFESH HA-CHAYAH))
(va-yar ELOHIM ki TOV)

(va-yevarekh ELOHIM otam                       ;; God blessed them
  (peru u-revu                                 ;; be fruitful, multiply
    :u-milu et HA-MAYIM
    :ve-ha-of yirev ba-ARETZ))

(va-yehi EREV va-yehi VOKER :day 5)

;; ══════════════════════════════════════════════════════════════
;; 1:24-31 — DAY 6: ANIMALS, HUMANS
;; ══════════════════════════════════════════════════════════════

(yomer ELOHIM
  (totze HA-ERETZ NEFESH CHAYAH le-minah
    (BEHEMAH)                                  ;; livestock
    (REMES)                                    ;; creeping things
    (CHAYTO-ERETZ le-minah)))                  ;; wild animals
(va-yehi KEN)

;; THE DIVINE COUNCIL — plural deliberation
(yomer ELOHIM
  (naaseh ADAM                                 ;; "Let US make human"
    :be-tzalmenu                               ;; in OUR image
    :ki-dmutenu))                              ;; according to OUR likeness

(va-yivra ELOHIM et HA-ADAM be-tzalmo          ;; God created human
  :be-tzelem ELOHIM bara oto                   ;; in image of God
  :zakhar u-nekevah bara otam)                 ;; male and female

(va-yevarekh ELOHIM otam
  (peru u-revu u-milu et HA-ERETZ              ;; be fruitful, fill earth
    :ve-kivshuha                               ;; subdue it
    :u-redu                                    ;; rule over
      (be-DAGAT HA-YAM)
      (u-ve-OF HA-SHAMAYIM)
      (u-ve-khol CHAYAH)))

(yomer ELOHIM
  (hineh natati lakhem                         ;; "Behold I give you"
    (et kol ESEV zorea ZERA)                   ;; seed-bearing plants
    (ve-et kol HA-ETZ asher-bo PRI-ETZ)        ;; fruit trees
    :lakhem yihyeh le-okhlah))                 ;; for food

(va-yar ELOHIM et kol asher ASAH               ;; God saw ALL made
  :ve-hineh TOV MEOD)                          ;; and behold: VERY GOOD

(va-yehi EREV va-yehi VOKER :day 6)

;; ══════════════════════════════════════════════════════════════
;; 2:1-3 — DAY 7: REST
;; ══════════════════════════════════════════════════════════════

(va-yekhulu HA-SHAMAYIM ve-HA-ARETZ            ;; heavens/earth finished
  :ve-khol TZEVAAM)                            ;; and all their host

(va-yekhal ELOHIM ba-YOM HA-SHEVII             ;; God finished day 7
  :melakhto asher ASAH)                        ;; work which he made

(va-yishbot ba-YOM HA-SHEVII                   ;; He rested day 7
  :mi-kol melakhto asher ASAH)

(va-yevarekh ELOHIM et YOM HA-SHEVII           ;; God blessed day 7
  (va-yekadesh oto)                            ;; sanctified it
  :ki vo SHAVAT mi-kol melakhto
  :asher bara ELOHIM la-asot)

;; ══════════════════════════════════════════════════════════════
;; END PROGRAM
;; ══════════════════════════════════════════════════════════════

(return UNIVERSE)
