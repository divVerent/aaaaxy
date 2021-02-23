sox \
  --combine merge \
  ../shepard_{00,04,07,00,04,07,00,07,04,00,07,04,00,04,07,00,04,07,00,07,04,00,07,04,00}.wav \
  ../loom.wav \
  delay $(seq 1 0.2 5.8) \
  remix - \
  remix - - \
  reverb 95 50 100 100 \
  fade q 0 7.8 1.5 \
  gain -n -1
