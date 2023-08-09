const Client = require('./client.js');

test('client sorts results from high-to-low', () => {
    const client = new Client('');
    const sorted = client._sortResults({
        scores: {
            'P1': 3,
            'P2': 1,
            'P3': 6,
        }
    });

    expect(sorted).toStrictEqual([
        {
            clientId: 'P3',
            score: 6
        },
        {
            clientId: 'P1',
            score: 3
        },
        {
            clientId: 'P2',
            score: 1
        }
    ]);
});