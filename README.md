# Projet RT805

Projet universitaire RT0805 : Initiation au protocole gRPC en Go.

## Prérequis

Assurez-vous d'installer les éléments suivants avant de commencer :
- **Go**
- **MongoDB**

## Configuration Docker

Pour créer un conteneur MongoDB, utilisez Docker Compose. 

- **MongoExpress** : Pour Accéder à l'interface à l'adresse `http://ip_du_conteneur:8081/`.
- **Identifiants** : Utilisez `admin` comme identifiant et `pass` comme mot de passe.

## Compilation

- **Serveur** : Pour lancer le serveur, allez dans le répertoire du serveur et exécutez `go run .`.
- **Client** : Pour lancer le client, allez dans le répertoire du client et exécutez `go run .`.

Le client envoie le fichier numéro 1 par défaut. Veillez à mettre à jour cette valeur si nécessaire.
