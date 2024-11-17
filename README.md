# Étude de Cas : Système de Commerce Électronique

## Présentation

Cette étude de cas a été réalisée pour mettre en pratique l'approche microservices, illustrant ainsi l'architecture et son fonctionnement. Ce projet représente une architecture microservices conçue pour gérer un système de commerce électronique.

## Auteur

- Léo Baleras
- Vincent Thome

## Lancer le Projet

Pour lancer ce projet, il vous suffit de suivre les étapes suivantes :

1. copier le fichier `.env.example` en `.env`
2. Exécutez le fichier `docker-compose.yml` pour démarrer les services.
3. Vous avez deux options pour interagir avec le système :
   - Allez dans le répertoire `/web-service` et exécutez `npm run dev` pour démarrer l'interface graphique.
     - pour se connecter, utilisez les identifiants suivants :
       - email : `user1@example.com`
       - mot de passe : `password`
   - Exécutez le script `sh run-test.sh` pour tester tous les endpoints des services en ligne de commande et afficher les résultats.

## Conclusion

Ce projet démontre l'efficacité de l'architecture microservices dans la gestion d'un système de commerce électronique, en permettant une scalabilité et une flexibilité accrues.
