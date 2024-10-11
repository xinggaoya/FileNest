import httpService from '@/api/http';

async function fetchData() {
  try {
    const data = await httpService.get('/data');
    console.log(data);
  } catch (error) {
    console.error('Error fetching data:', error);
  }
}

async function postData() {
  try {
    const response = await httpService.post('/data', { key: 'value' });
    console.log(response);
  } catch (error) {
    console.error('Error posting data:', error);
  }
}
