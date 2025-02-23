# Distribution Management System

This project was developed as part of a machine task for a company interview process. It implements a distribution management system with features for managing distributors and their permissions across different regions.

## Core Features

- **Distributor Management**
  - Add new distributors
  - Remove existing distributors
  - Create sub-distributor relationships
  - Region-based validation using cities.csv data

- **Permission Control**
  - Allow distribution rights
  - Revoke distribution rights
  - Check distribution permission status

- **Validation System**
  - Region validation against cities.csv database
  - Distributor existence verification

## API Endpoints

### Distributor Management

#### 1. Add Distributor
- **Endpoint**: `POST /distributor`
- **Description**: Register a new distributor in the system
- **Request Body**:
```json
{
    "distributor":"distributer_name"
}
  ```
- **Success Response**: 201 Created

#### 2. Remove Distributor
- **Endpoint**: `DELETE /distributor/:distributor`
- **Description**: Remove an existing distributor from the system
- **Path Parameter**: `distributor` - Name of the distributor
- **Success Response**: 200 OK

#### 3. Add Sub-Distributor
- **Endpoint**: `POST /distributor/add-sub`
- **Description**: Create a hierarchical relationship between distributors
- **Request Body**:
  ```json
{
    "parent_distributor":"distributer_name",
    "sub_distributor":"new_sub_distributer_name""
}
  ```
- **Success Response**: 201 Created

### Permission Management

#### 1. Check Distribution Permission
- **Endpoint**: `GET /permission/check`
- **Description**: Verify if a distributor has permission for a specific region
- **Query Parameters**: 
  - `distributor`: Distributor name
  - `region`: Region to check
- **Success Response**: 200 OK with permission status

#### 2. Allow Distribution
- **Endpoint**: `POST /permission/allow`
- **Description**: Grant distribution rights to a distributor
- **Request Body**:
  ```json
  {
    "distributor": "distributor_name",
    "region": "region_name" (Eg: "KLRAI-TN-IN")
  }
  ```
- **Success Response**: 200 OK

#### 3. Disallow Distribution
- **Endpoint**: `POST /permission/disallow`
- **Description**: Revoke distribution rights from a distributor
- **Request Body**:
  ```json
  {
    "distributor": "distributor_name",
    "region": "region_name"
  }
  ```
- **Success Response**: 200 OK

## Technical Implementation

### Implementation Details
- **Clean Architecture Pattern**
  - Separation of concerns
    - Route handlers in `internal/handler`
    - Distribution logics and saving data handled in `internal/data`
  - Easy to test and maintain
- **Validation System**
  - Region validation using CSV data
  - Distributor existence checks
  - Permission hierarchy validation
- **Concurrency Handling**
  - Prevention of race conditions during concurrent access using read and write locks (sync.RWMutex)

### Key Components
1. **Route Handlers** (`internal/handler`)
   - Handle HTTP requests and responses
   - Input validation and sanitization
   - Error handling and response formatting

2. **Business Logic & Data Management** (`internal/data`)
   - Distributor management logic
   - Permission validation and inheritance
   - Region validation against CSV data
   - Thread-safe data operations using RWMutex
   - Optimized read/write locking for better performance

## Technical Notes

- The system performs validation against a predefined list of cities/regions from cities.csv
- Distributor authentication is simplified (no session management) with distributer name passed in request body
- All operations include validation for distributor existence and permission checks
- The system maintains hierarchical relationships between distributors and sub-distributors

## Future Improvements

Potential enhancements that could be added:
1. Proper authentication and session management
2. Database persistence for distributor and permission data (restricted by assignment)
3. Caching layer for frequently accessed data (restricted by assignment)
4. More detailed logging and monitoring
5. Rate limiting for API endpoints
6. Testing(unit and integrated)

## Proposed Model Improvements

1. **Enhanced Distribution Model**
   - Create separate entities for distributors
   - Enable multi-directional distribution relationships