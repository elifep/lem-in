# Lem-in

## Project Description
This project simulates the behavior of ants in an ant farm. The objective is to find the optimal paths for a given number of ants to travel from a starting point to a destination without causing traffic jams.

The simulation aims to:
- Use Depth-First Search (DFS) to explore possible paths.
- Avoid redundant or intersecting paths.
- Calculate the best combinations of paths based on the number of ants and path lengths.
- Assign ants to tunnels efficiently, ensuring that only one ant is in a tunnel at a time.

## Approach
The algorithm follows these steps:

1. **Path Finding with DFS**: DFS is used to identify all possible paths from the start to the destination.
2. **Path Optimization**: Paths found in the first step are compared to avoid redundancy, retaining only the shortest, non-overlapping paths.
3. **Combination Selection**: For different numbers of ants, the algorithm calculates and selects path combinations with the shortest total length.
4. **Ant Assignment**: Based on the selected combination, ants are assigned to tunnels, with only one ant in each tunnel at any time.
5. **Path Traversal**: Ants travel through their assigned paths to reach the destination.

## Usage

### Requirements
- Go programming language installed on your system.

### Running the Simulation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/elifep/lem-in.git
   cd lem-in
2. **Run the Simulation:** Run the simulation using the following command, replacing [filename.txt] with the name of your input file:
3. 
go run . filename.txt
4. **Example Outputs**
Here’s an example of how the output may look based on different files:

**EXAMPLE03**
go run . example03.txt

All paths from 0 to 5 :
Path: 0 -> 1 -> 4 -> 5
Path: 0 -> 2 -> 4 -> 5

Maximum non-overlapping paths:
0 -> 1 -> 4 -> 5
Tur sayısı: 6
L1-1
L1-4 L2-1
L1-5 L2-4 L3-1
L2-5 L3-4 L4-1
L3-5 L4-4
L4-5

Execution time: 0.00020972 seconds

4. **Testing with 01 Edu Cases**
For testing with predefined cases from 01 Edu, visit [this link](https://github.com/01-edu/public/tree/master/subjects/lem-in/audit). Download the test cases and run them using the command above.

5. **Project Structure**
```text
lem-in/
├── main.go                  # Main Go file to run the simulation
├── README.md                # Project documentation
├── exampleXX.txt            # Example input files for different cases
└── go.mod                   # Go module file
```
6. ## Author
This project was created in May 2024 by Elif Ep.

7. ## Contributing
Contributions are welcome! If you'd like to contribute to this project, please follow these steps:

1. **Fork the Repository**: Click on the "Fork" button at the top right of this page to create a copy of this repository under your own GitHub account.

2. **Clone the Forked Repository**:
   ```bash
   git clone https://github.com/elifep/lem-in.git
   cd lem-in
   
   Create a New Branch:
   ```bash
   git checkout -b feature-or-bugfix-name

   Make Changes: Develop your feature or fix the bug, and make sure your code is clean and tested.

   Commit Changes:
   ```bash
   git commit -m "Description of your changes"

   Push to Your Forked Repository:
   ```bash
   git push origin feature-or-bugfix-name
   
   Create a Pull Request: Go to the original repository and open a pull request with a detailed description of your changes.
