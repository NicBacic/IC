# Projeto de iniciação científica desenvolvido por Nicolas Bacic e Daniel Cordeiro.

Nesse projeto estudamos o problema conhecido na literatura como MOSP, um problema de escalonamento que envolve múltiplas organizações que participam de uma mesma plataforma cooperativa. 7

Organizações, nesse contexto, é qualquer entidade que possua um conjunto de tarefas que devem ser processadas e que envolvam algum problema de escalonamento

Estudamos uma variante do problema MOSP onde a função objetivo das organizações é minimizar o consumo de energia de seus computadores. 

Nesse problema temos um conjunto de tarefas definidas por release date, deadline e volume de processamento (o quanto de "trabalho" o processador deve exercer para executar a tarefa inteiramente).

O consumo de energia está diretamente associado a velocidade de execução das tarefas. Quanto menor a velocidade, menos energia é consumida. Portanto basta determinar qual é a menor velocidade que podemos executar as tarefas de tal forma que elas não atrasem e sejam executadas por completo.

Os algoritmos nesse projeto foram inteiramente desenvolvidos na linguagem Go. O algoritmo ILBA parte de uma heurística de balancear as organizações mais carregadas a cada iteração. Nele, a função objetivo das organizações é o makespan

O algoritmo YDS é um algoritmo de escalonamento de tarefas para obter o consumo de energia mínimo. Ele funciona por meio de cálculos repetitivos do intervalo de execução das tarefas.

Por último temos o algoritmo MOSPEnergy, o algoritmo de escalonamento de tarefas em múltiplas organizações. Ele utiliza a heurística ILBA combinado com o algoritmo YDS.
