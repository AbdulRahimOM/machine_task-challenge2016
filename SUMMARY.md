# 🏢 Distribution Management System

This project was developed as part of a machine task for a company interview process. It implements a distribution management system with features for managing distributors and their permissions across different regions.

## 🎯 Core Features

📦 **Distributor Management**
  - Add new distributors
  - Remove existing distributors
  - List all distributors

🔑 **Permission Management**
  - Allow distribution rights over a region
  - Disallow distribution rights over a region
  - Check permission status over specific regions (Responses: FULLY_ALLOWED/PARTIALLY_ALLOWED/FULLY_DENIED)
  - View distributor-specific permissions (As text(contract) or JSON)
  - Contract-based permission management

🌍 **Region Management**
  - Hierarchical region structure (Country → Province → City)
  - Region validation against cities.csv database
  - Region code format: "CITYCODE-PROVINCECODE-COUNTRYCODE"

## 🚀 How to use

### Prerequisites  
- Go 1.23 or higher  
- Git  

### Installation
1. Clone the repository
```bash
git clone https://github.com/AbdulRahimOM/machine_task-challenge2016.git
cd machine_task-challenge2016
```

2. Set up environment variables
```bash
touch .env
echo PORT="4010" >> .env # Or any other port number
```

3. Build the project
```bash
make build
```

4. Run the server
```bash
./bin/app
```

The server will start on `localhost:4010` (or the port specified in the .env file).

### Region Format
• Countries: 2-letter code (e.g., "IN", "US")
• Provinces: 2-letter code + country (e.g., "TN-IN")
• Cities: City code + province + country (e.g., "CENAI-TN-IN")

## 🛠️ API Endpoints

### 📦 Distributor Management

#### 1. Add Distributor
- **Endpoint**: `POST /distributor`
- **Description**: Register a new distributor in the system
- **Request Body**:
  ```json
  {
    "distributor": "distributor_name"
  }
  ```
- **Success Response**: 201 Created

#### 2. Remove Distributor
- **Endpoint**: `DELETE /distributor/:distributor`
- **Description**: Remove an existing distributor from the system
- **Path Parameter**: `distributor` - Name of the distributor
- **Success Response**: 200 OK

#### 3. Get Distributors
- **Endpoint**: `GET /distributor`
- **Description**: Retrieve list of all distributors
- **Success Response**: 200 OK with distributors list

### 🔑 Permission Management

#### 1. Check Distribution Permission
- **Endpoint**: `GET /permission/check`
- **Description**: Verify distribution permission status for a region
- **Query Parameters**: 
  - `distributor`: Distributor name
  - `region`: Region to check
- **Success Response**: 200 OK with permission status

#### 2. Allow Distribution
- **Endpoint**: `POST /permission/allow`
- **Description**: Grant distribution rights for a region
- **Request Body**:
  ```json
  {
    "distributor": "distributor_name",
    "region": "region_name" // Example: "KLRAI-TN-IN"
  }
  ```
- **Success Response**: 200 OK

#### 3. Apply Contract
- **Endpoint**: `POST /permission/contract`
- **Description**: Apply distribution contract with permissions
- **Success Response**: 200 OK

#### 4. Disallow Distribution
- **Endpoint**: `POST /permission/disallow`
- **Description**: Revoke distribution rights
- **Request Body**:
  ```json
  {
    "distributor": "distributor_name",
    "region": "region_name"
  }
  ```
- **Success Response**: 200 OK

#### 5. Get Distributor Permissions
- **Endpoint**: `GET /permission/:distributor`
- **Description**: Retrieve all permissions for a distributor in either JSON or contract text format
- **Path Parameter**: `distributor` - Name of the distributor
- **Query Parameter**: `type` - Response format type ("json" or "text")
  - `json`: Returns structured JSON format with permissions
  - `text`: Returns formatted contract-like text representation
- **Success Response**: 200 OK with permissions in requested format
- **Response Examples**:
  - Text format (`type=text`):
    ```text
    Permissions for DISTRIBUTOR1
    INCLUDE: IN
    INCLUDE: US
    INCLUDE: ONATI-SS-ES
    EXCLUDE: KA-IN
    EXCLUDE: CENAI-TN-IN
    ```
  - JSON format (`type=json`):
    ```json
    {
        "status": true,
        "resp_code": "SUCCESS",
        "data": {
            "Distributor": "DISTRIBUTOR1",
            "Included": [
                "IN",
                "US",
                "ONATI-SS-ES"
            ],
            "Excluded": [
                "KA-IN",
                "CENAI-TN-IN"
            ]
        }
    }
    ```

### 🌍 Region Management

#### 1. Get Countries
- **Endpoint**: `GET /regions/countries`
- **Description**: Get list of available countries
- **Success Response**: 200 OK with countries list

#### 2. Get Provinces
- **Endpoint**: `GET /regions/provinces/:countryCode`
- **Description**: Get provinces in a country
- **Path Parameter**: `countryCode`
- **Success Response**: 200 OK with provinces list

#### 3. Get Cities
- **Endpoint**: `GET /regions/cities/:countryCode/:provinceCode`
- **Description**: Get cities in a province
- **Path Parameters**: 
  - `countryCode`
  - `provinceCode`
- **Success Response**: 200 OK with cities list

## 🏗️ Technical Implementation

### 🎨 Architecture  
- **Clean Architecture Pattern**  
  - Separation of concerns with handlers and business logic  
  - RESTful API design  
  - Modular component structure  

### 🔧 Key Components  
1. **Route Handlers** (`internal/handler`)  
   - HTTP request handling  
   - Input validation  
   - Response formatting  
   - Error handling  

2. **Data Management**  
   - In-memory data storage  
   - Thread-safe operations using `sync.RWMutex`  
   - CSV-based region validation  
   - Contract validation and processing  
  
### ⚙️ Technical Features
- Region validation against cities.csv
- Concurrent access handling with sync.RWMutex
- Hierarchical permission system
- Contract-based permission management
- Region-based distribution control

## 📝 Technical Notes
- Thread-safe operations using read-write mutex locks
- CSV-based region validation
- Hierarchical region structure validation
- Contract template validation

## 🚀 Potential Improvements (if assignment is flexible)

⏳ **Contract-expiry**  
   - When contract expires, cascade expiration to all dependent sub-contracts  
   - Inheritance on contract, and not on permission  
    ↳ This would be more matching to the real-world scenario, where permissions are time-based and amendable contracts  

🔐 **Distributor Self-Service Portal**  
   - Implement secure authentication system  
   - Enable distributors to manage their own sub-contracts  
    ↳ Create and modify sub-contracts within their permitted scope  
    ↳ Monitor contract status and expiration dates  
    ↳ View inheritance chain and dependencies  
    ↳ Notify them when contract expires  