# TP2: Diving deeper with Neo4j

## Sommaire

- [Lecture du fichier](#lecture-du-fichier)
- [Importation des données](#importation-des-données)
- [Résultat](#résultat)
- [Utilisation](#utilisation)
  - [Paramètres](#paramètres)
  - [Native](#native)
  - [Docker](#docker)
  - [Kubernetes](#kubernetes)

## Lecture du fichier

Le programme prend en charge les fichiers à partir d'une source locale ou d'un lien vers une archive 7zip. Dans le cas d'un lien distant, le programme télécharge l'archive et la décompresse.

La lecture du fichier s'effectue en streaming, c'est-à-dire que le fichier est lu par blocs de 512 octets, et le programme tente de créer des objets JSON valides correspondant aux articles.

Il faut noter que le fichier JSON initial présentait des problèmes de format, notamment des valeurs telles que `NumberInt(1)` dans certaines propriétés comme "year", "venue.type" ou "n_citation". Pour résoudre cette anomalie, j'ai mis en place un lecteur qui lit 19 octets supplémentaires atteindre la longueur de "NumberInt()". Ensuite, une vérification avec une expression régulière (regex) est effectuée pour supprimer ce texte et ne laisser que le chiffre et ainsi rendre la valeur conforme à un objet JSON standard.

En plus de ce problème, il manquait également l'identifiant à certains auteurs. Pour remédier à cela, je leur ai rajouté un identifiant en prenant l'ID de l'article et en y ajoutant un nombre en hexadécimal à la fin.

## Importation des données

Des contraintes ont été établies sur les IDs des articles et des auteurs afin de garantir leur unicité et d'améliorer la rapidité des recherches grâce aux index.

Le programme fonctionne en parallèle avec deux threads. Tout d'abord, un thread est responsable de la lecture en continu du fichier et crée des lots de 2'000 articles. Une fois le lot créé, il le place dans une file d'attente qui sera traitée par le deuxième thread. 
Chaque lot d'articles est ensuite traité dans une transaction qui effectue les opérations suivantes:
1. Création du noeud d'un article s'il n'existe pas, et on défini son titre.
2. Création du noeud pour chaque auteur de l'article, et si on viens de le créer, on défini son nom.
3. On créer la relation "AUTHORED" avec l'auteur.
4. Création de l'article cité s'il n'existe pas.
5. Enfin, on créer la relation "CITE" vers cette article.

## Résultat

{"team"="BenammAdvDaBa23", "N"="12115473", "RAM_MB"="3000", "seconds"="21393"}

## Utilisation

### Paramètres

Le programme prend en charge plusieurs paramètres. Voici une liste des options disponibles :

| Paramètre | Variable d'environnement | Valeur par défaut | Description |
|---|---|---|---|
| host | DATABASE_HOST | localhost:7687 | Hôte de la base de données |
| user | DATABASE_USER | neo4j | Nom d'utilisateur de la base de données |
| pass | DATABASE_PASS | aztec-peace-linear-laura-gregory-4537 | Mot de passe de la base de données |
| file | FILE_PATH | biggertest.json | Chemin du fichier |
| download-file | DOWNLOAD_FILE |  | Si défini, URL de l'archive 7z à télécharger |
| size | BATCH_SIZE | 2000 | Taille des lots |

Utilisation des paramètres: 
```bash
go run . --parameter value 
```

Utilisation des variables d'environnement: 
```bash
export VARIABLE=value
```

### Native

1. [Installez Neo4J](https://neo4j.com/docs/operations-manual/current/installation/) et configurez vos identifiants.
2. [Installez Golang](https://go.dev/doc/install). 
3. Exécutez le programme en spécifiant vos identifiants avec la commande suivante :
```bash 
go run . --user neo4j --pass test --download-file https://originalstatic.aminer.cn/misc/dblp.v13.7z
```

### Docker

Vous avez deux possibilités : soit exécuter les deux conteneurs individuellement, soit les déployer simultanément à l'aide de Docker Compose.

#### Utilisation de Docker Compose

```bash
docker compose -f docker/docker-compose.yaml up
```

#### Ou de manière individuelle

1. Construire l'image Docker de l'indexeur:
```bash
docker build -t indexer -f docker/Dockerfile .
```

2. Lancer la base de données Neo4j:
```bash
docker run -p 7474:7474 -p 7687:7687 neo4j:latest --name database
```

3. Lancer l'indexeur:
```bash
docker run indexer -e DATABASE_USER=neo4j -e DATABASE_PASS=test -e DOWNLOAD_FILE=https://originalstatic.aminer.cn/misc/dblp.v13.7z --name indexer
```

Ces étapes vous permettront de déployer et exécuter les conteneurs nécessaires à l'indexation des données dans Neo4j. L'option Docker Compose simplifie ce processus en orchestrant le démarrage simultané des deux conteneurs.

### Kubernetes

1. Définir la configuration KubeConfig du cluster :
```bash
export KUBECONFIG=local.yaml
```

2. Déployer la base de données Neo4j:
```bash
kubectl apply -f kubernetes/database-deployment.yaml --namespace your-namespace
```

3. Déployer l'indexeur*:
```bash
kubectl apply -f kubernetes/indexer-deployment.yaml --namespace your-namespace
```
* Ce déploiement Kubernetes utilise une image déployée sur github : si vous voulez modifier le code vous devez modifier le fichier et le changer avec votre propre image