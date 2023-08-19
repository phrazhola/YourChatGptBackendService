# Build Your ChatGPT

Welcome! 

YourChatGPT is a versatile AI chatbot solution that harnesses the capabilities of the ChatGPT API. Crafted to simplify your journey, it enables you to create a tailored ChatGPT clone effortlessly.

Feel free to visit our website to explore a comprehensive introduction and experience a demo of the sample Chatbot in action: [Build Your Own ChatGPT]()

**Frontend package**: https://github.com/phrazhola/YourChatGptReactApp

## Highlights
-   **Structured Code:**  The codebase is well-organized, making it easy to understand and customize for different styles and use cases.
-   **Extensible API:**  The API is designed with extensibility in mind, offering flexibility for future enhancements.
-   **Scalability:**  The architecture supports scalability, with a clear separation between the frontend UI and backend API.
-   **Plug-and-Play:**  The solution is ready to use, requiring only API key setup.
-   **Security First:**  The system prioritizes security by ensuring your OpenAI API key remains secure on the client side.

# Setup

There are three things you need to setup :

- OpenAI API key
- Database
- Middleware

## OpenAI API Key
- Go to https://platform.openai.com/docs/api-reference/authentication and there is guide of how to get your API key.
- Keep your API key secret and declare it as environment variable in your deployment environment and access it inside your application code, or specify it directly in your code (`/clients/openai_client.go`) (NOT recommended).

## Database
- We leverage Azure Cosmos DB as our NoSQL database solution. However, you have the flexibility to opt for your preferred alternative and adapt the code accordingly.
- If you decide to continue using Azure Cosmos DB, there's a free plan available that you can take advantage of. For comprehensive guidance, you can refer to the comprehensive documentation provided at: [Azure Cosmos DB for NoSQL documentation](https://learn.microsoft.com/en-us/azure/cosmos-db/nosql/).
- To set up your database within the Azure portal, you can follow this step-by-step tutorial: [Quickstart: Create an Azure Cosmos DB account, database, container, and items from the Azure portal](https://learn.microsoft.com/en-us/azure/cosmos-db/nosql/quickstart-portal).
- Feel free to choose the database solution that suits your needs best, and modify the code as necessary. If you opt for Azure Cosmos DB, the documentation provided will guide you through the setup and integration process.
- You will need to clarify the connectionString, databaseID and containerID in `/database/util.go`.


## Middleware for the Gin web framework
- Pretty much you will just need to add your frontend application endpoint to the files under `/middleware/`.

