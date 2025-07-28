```mermaid
graph TD
    User[End User] --> GlobalRouter[Azure Traffic Manager or Front Door]

    subgraph Azure US East Region
        direction LR
        GlobalRouter --> VNet_AKS_East[VNet East AKS]
        subgraph VNet_AKS_East
            direction TB
            Firewall_AKS_East[Azure Firewall or NSGs]

            subgraph AKS_Cluster_East [AKS Cluster East]
                direction TB
                Firewall_AKS_East --> Ingress_East[K8s Ingress Controller]

                subgraph ns_alpha_east [Namespace App Alpha]
                    style ns_alpha_east fill:#f1f8e9,stroke:#558b2f
                    Ingress_East -- /alpha --> Alpha_SvcFE_East[Service FE Alpha]
                    Alpha_SvcFE_East --> Alpha_DeployFE_East[Deployment FE Alpha]
                    Alpha_DeployFE_East --> Alpha_SvcBE_East[Service BE Alpha]
                    Alpha_SvcBE_East --> Alpha_DeployBE_East[Deployment BE Alpha]
                    Alpha_DeployBE_East -- Private Endpoint --> Alpha_DB_East[Azure SQL or Postgres Alpha DB]
                end

                subgraph ns_bravo_east [Namespace App Bravo]
                     style ns_bravo_east fill:#f1f8e9,stroke:#558b2f
                     Ingress_East -- /bravo --> Bravo_SvcFE_East[Service FE Bravo]
                     Bravo_SvcFE_East --> Bravo_DeployFE_East[Deployment FE Bravo]
                     Bravo_DeployFE_East --> Bravo_SvcBE_East[Service BE Bravo]
                     Bravo_SvcBE_East --> Bravo_DeployBE_East[Deployment BE Bravo]
                     Bravo_DeployBE_East -- Private Endpoint --> Bravo_DB_East[Azure SQL or Postgres Bravo DB]
                 end
            end
            Alpha_DB_East -- VNet Integration --> VNet_AKS_East
            Bravo_DB_East -- VNet Integration --> VNet_AKS_East
        end
    end

    subgraph Azure US West Region
        direction LR
        GlobalRouter --> VNet_AKS_West[VNet West AKS]
         subgraph VNet_AKS_West
            direction TB
            Firewall_AKS_West[Azure Firewall or NSGs]

            subgraph AKS_Cluster_West [AKS Cluster West]
                direction TB
                Firewall_AKS_West --> Ingress_West[K8s Ingress Controller]

                subgraph ns_alpha_west [Namespace App Alpha]
                    style ns_alpha_west fill:#f1f8e9,stroke:#558b2f
                    Ingress_West -- /alpha --> Alpha_SvcFE_West[Service FE Alpha]
                    Alpha_SvcFE_West --> Alpha_DeployFE_West[Deployment FE Alpha]
                    Alpha_DeployFE_West --> Alpha_SvcBE_West[Service BE Alpha]
                    Alpha_SvcBE_West --> Alpha_DeployBE_West[Deployment BE Alpha]
                    Alpha_DeployBE_West -- Private Endpoint --> Alpha_DB_West[Azure SQL or Postgres Alpha DB]
                end

                subgraph ns_bravo_west [Namespace App Bravo]
                     style ns_bravo_west fill:#f1f8e9,stroke:#558b2f
                     Ingress_West -- /bravo --> Bravo_SvcFE_West[Service FE Bravo]
                     Bravo_SvcFE_West --> Bravo_DeployFE_West[Deployment FE Bravo]
                     Bravo_DeployFE_West --> Bravo_SvcBE_West[Service BE Bravo]
                     Bravo_SvcBE_West --> Bravo_DeployBE_West[Deployment BE Bravo]
                     Bravo_DeployBE_West -- Private Endpoint --> Bravo_DB_West[Azure SQL or Postgres Bravo DB]
                 end
            end
            Alpha_DB_West -- VNet Integration --> VNet_AKS_West
            Bravo_DB_West -- VNet Integration --> VNet_AKS_West
        end
    end

    note right of VNet_AKS_East: VNet and Firewall and Private Endpoints and K8s Network Policies provide HIPAA SOC2 boundary
    note right of VNet_AKS_West: VNet and Firewall and Private Endpoints and K8s Network Policies provide HIPAA SOC2 boundary

    %% Basic Styling (No symbols in text, no underline)
    classDef vnet fill:#e3f2fd,stroke:#1e88e5,color:black;
    classDef fw fill:#ffebee,stroke:#c62828,color:black;
    classDef aks fill:#eceff1,stroke:#37474f,color:black;
    classDef k8s_compute fill:#e0f7fa,stroke:#006064,color:black;
    classDef k8s_svc fill:#e0f2f1,stroke:#004d40,color:black;
    classDef k8s_ingress fill:#ede7f6,stroke:#311b92,color:black;
    classDef db fill:#f3e5f5,stroke:#6a1b9a,color:black;
    classDef global fill:#e0f7fa,stroke:#006064,color:black;

    class VNet_AKS_East,VNet_AKS_West vnet;
    class Firewall_AKS_East,Firewall_AKS_West fw;
    class AKS_Cluster_East,AKS_Cluster_West aks;
    class Alpha_DeployFE_East,Alpha_DeployBE_East,Bravo_DeployFE_East,Bravo_DeployBE_East,Alpha_DeployFE_West,Alpha_DeployBE_West,Bravo_DeployFE_West,Bravo_DeployBE_West k8s_compute;
    class Alpha_SvcFE_East,Alpha_SvcBE_East,Bravo_SvcFE_East,Bravo_SvcBE_East,Alpha_SvcFE_West,Alpha_SvcBE_West,Bravo_SvcFE_West,Bravo_SvcBE_West k8s_svc;
    class Ingress_East,Ingress_West k8s_ingress;
    class Alpha_DB_East,Bravo_DB_East,Alpha_DB_West,Bravo_DB_West db;
    class GlobalRouter global;
    class User fill:#f8bbd0,stroke:#880e4f,color:black;
```