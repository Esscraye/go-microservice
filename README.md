# Microservices Architecture Case Study

## Lancement du projet

Pour lancer le projet il vous faut avoir docker et docker-compose installé sur votre machine.

ensuite vous pouvez lancer les microservices avec la commande suivante:

```bash
docker-compose up -d --build
```

Un script shel (run-test.sh) est disponible pour tester les services.

```bash
./run-test.sh
```

Et pour arrêter et supprimer complètement les docker containers (attention si vous avez d'autres volumes docker la commande suivante les supprimera aussi):

```bash
docker compose down --remove-orphans -v --rmi all && docker system prune -af && docker volume prune -af
```

## Group Members

- Member 1
- Member 2
- Member 3

## Project Description

This project demonstrates the implementation of a microservices architecture using the Go programming language. The application consists of multiple services that communicate via APIs.

## Application de Gestion de Commandes en Ligne

Cette application permettrait de gérer les commandes en ligne pour un magasin. Elle pourrait inclure plusieurs microservices, chacun responsable d'une fonctionnalité spécifique.

## Structure des Microservices

1. Service de Gestion des Produits (Product Service)

- CRUD pour les produits.
- Recherche de produits par catégorie, nom, etc.

2. Service de Gestion des Utilisateurs (User Service)

- CRUD pour les utilisateurs.
- Authentification et autorisation.

3. Service de Gestion des Commandes (Order Service)

- Création et gestion des commandes.
- Suivi des commandes.

4. Service de Paiement (Payment Service)

- Traitement des paiements.
- Gestion des transactions.

5. Service de Notification (Notification Service)

- Envoi de notifications par email ou SMS pour les confirmations de commande et les mises à jour.