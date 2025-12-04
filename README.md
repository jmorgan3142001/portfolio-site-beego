# jake-morgan.dev

The source code for my personal portfolio and digital dashboard.

This repository houses my personal portfolio application. I developed this project as a solo contributor to demonstrate my ability to architect and deploy modern distributed systems from the ground up. It serves as a practical showcase of my full stack engineering skills in a live production environment.

## Technical Architecture

I designed the system to be stateless and scalable by leveraging cloud-native technologies.

* **Backend Framework:** I chose the Beego framework with Go (Golang) to utilize strict typing and high concurrency performance. The MVC structure allows for organized business logic and maintainable code.
* **Cloud Infrastructure:** The application is containerized and hosted on Google Cloud Cloud Run. This serverless approach ensures automatic scaling and high availability with minimal operational overhead.
* **Data Persistence:** I integrated Supabase as the primary database layer. This handles relational data requirements securely and efficiently while interacting seamlessly with the Go backend.
* **External Integrations:** The backend orchestrates requests to multiple third party APIs. This aggregates data from various external sources to present a unified and dynamic view of my professional activity.

## Live Deployment

The application is currently live and serving traffic. You can access the main dashboard here:

**[🔗 Launch Terminal (jake-morgan.dev)](https://jake-morgan.dev/)**