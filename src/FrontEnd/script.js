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
      // Display response in the results div
      const resultsDiv = document.getElementById('results');
      resultsDiv.innerHTML = `
        <p>Solution: ${data.solution}</p>
        <p>Articles Checked: ${data.articlesChecked}</p>
        <p>Path Length: ${data.pathLength}</p>
        <p>Time Taken: ${data.timeTaken}</p>
      `;
    })
    .catch(error => {
      console.error('Error:', error);
    });
  });
  