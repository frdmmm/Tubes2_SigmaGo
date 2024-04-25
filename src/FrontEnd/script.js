document.getElementById('searchForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const start = document.getElementById('start').value;
    const end = document.getElementById('end').value;
    const algorithm = document.getElementById('algorithm').value;
  
    // Construct form data
    const formData = new FormData();
    formData.append('start', start);
    formData.append('end', end);
    formData.append('algo', algorithm);
  
    // Make a POST request to the backend
    fetch('http://localhost:8080/solve', {
      method: 'POST',
      body: formData
    })
    .then(response => {
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      return response.json();
    }) // Parse response as JSON
    .then(data => {
      console.log(data);
      // Display response in the results div
      const resultsDiv = document.getElementById('results');
      resultsDiv.innerHTML = `
        <p>Solution: ${data.solution}</p>
        <p>Articles Checked: ${data.articlesChecked}</p>
        <p>Path Length: ${data.pathLength}</p>
        <p>Time Taken: ${data.timeTaken}</p>
      `;

      // Draw graph
      const graphContainer = document.getElementById('graph-container');
      graphContainer.innerHTML = ''; // Clear previous graph

      const solution = data.solution;
      const nodePositions = {};
      const nodeDistance = 100; // Adjust node distance as needed

      solution.forEach((node, index) => {
          const leftPosition = index * nodeDistance + 50; // Adjust starting position as needed
          nodePositions[node] = { left: leftPosition };
          
          // Draw node
          const nodeElement = document.createElement('div');
          nodeElement.textContent = node;
          nodeElement.className = 'node';
          nodeElement.style.left = `${leftPosition}px`;
          graphContainer.appendChild(nodeElement);
          
          // Draw edge
          if (index > 0) {
              const edgeElement = document.createElement('div');
              edgeElement.className = 'edge';
              edgeElement.style.left = `${leftPosition - nodeDistance / 2}px`; // Adjust edge position as needed
              graphContainer.appendChild(edgeElement);
          }
      });

      // Set edge height to match node height
      const edgeElements = document.querySelectorAll('.edge');
      edgeElements.forEach(edge => {
          edge.style.height = `${graphContainer.querySelector('.node').offsetHeight}px`;
      });
    })
    .catch(error => {
      console.error('Error:', error);
    });
});
