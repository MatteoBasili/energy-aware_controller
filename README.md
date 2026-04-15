# Energy-Aware Controller

Questo repository contiene l'implementazione del controllore energy-aware, basato sull’algoritmo **EAMO (Energy-Aware Microservices Orchestrator)**, per la tesi di laurea magistrale in ingegneria informatica.

Il sistema distribuisce dinamicamente il traffico tra varianti di uno stesso microservizio, ottimizzando il trade-off tra qualità del servizio e consumo energetico. Le decisioni sono basate su metriche raccolte da Prometheus (Istio e Kepler) e applicate tramite routing pesato nella service mesh.