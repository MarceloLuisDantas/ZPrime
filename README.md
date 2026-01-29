# ZPrime
Projeto ZPrime. Um emulador um computador retro fictício, inspirado em sistemas clássicos como NES, Atari 2600, Master System... 

## Objetivo
ZPrime possui como objetivo primario, prover uma experiencia proxima ao desenvolvimento em um sistema primitivo, sem a necessidade de lidar com conceitos mais complexos como V/H Blank, memory bank switch e sincronização de PPU, e para ser uma porta de entrada para entender melhor o funcionamento de linguagens Assemble, sem precisar entender sobre formato de executavel, layout de memoria ou por menores de cada sistema moderno. Apenas escreva seu programa, e execute ele de forma simples e rapida.

## Arquitetura

### Processador
ZPrime possui um processador 16 Bits capaz de lidar diretamente com inteiros de 8 e 16 com e sem sinal. Para evitar problemas com sincronização com o renderização, o processador não é limitado a alguma frequencia especifica, e o video é apenas renderizado a vontade do programador com o uso de syscalls, no entanto, o video é limitado a 60 FPS. O motivo para limitar o FPS a 60, é primeriamente evitar complicações com comuns em sistemas antigos sobre sincronização da logica do programa, com a renderização, e a dependencia no V/H Blank. 

### Memoria
Toda a memoria é implementada em arquitetura multi ship, com RAM, VRAM e ROM tendo cada um chip de memoria dedicado. O motivo para utilização de multi ship, ao invês de um unico chip (que é o mais comum em sistemas similares), é evitar problemas com mapeamento de memoria e bamk switch, invês de precisar lembrar que a RAM de uso livre inicia em 0x9000 e vai ate 0xF000 ou outros ranges para outro componentes, cada aspecto do sistema possui seu proprio chip que começa em 0x0000 e vai ate 0xFFFF. Unica exceção, é os valore de entrada que estão localizados nos 16 primeiros endereços da RAM.

### Display
ZPrime possui um display ASCII 60x60, com capacidade para 8 cores em texto e background.

## Aprenda
Todas as instruções estão documentadas em [Documentação](./Documentação.md)
