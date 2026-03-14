import apiClient from '../client';
import {getResources, getSections, getItems} from '../content';

jest.mock('../client', () => ({
  __esModule: true,
  default: {
    get: jest.fn(),
  },
}));

const mockedGet = apiClient.get as jest.Mock;

describe('Content API', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('getResources', () => {
    it('calls the correct endpoint with default params', async () => {
      const mockResponse = {
        data: {
          data: [],
          page: 1,
          page_size: 20,
          total_count: 0,
          total_pages: 0,
        },
      };
      mockedGet.mockResolvedValue(mockResponse);

      const result = await getResources();

      expect(mockedGet).toHaveBeenCalledWith('/api/resources', {
        params: {page: undefined, page_size: undefined},
      });
      expect(result).toEqual(mockResponse.data);
    });

    it('passes pagination params correctly', async () => {
      const mockResponse = {
        data: {
          data: [{id: 1, title: 'Test Resource'}],
          page: 2,
          page_size: 10,
          total_count: 15,
          total_pages: 2,
        },
      };
      mockedGet.mockResolvedValue(mockResponse);

      const result = await getResources({page: 2, page_size: 10});

      expect(mockedGet).toHaveBeenCalledWith('/api/resources', {
        params: {page: 2, page_size: 10},
      });
      expect(result).toEqual(mockResponse.data);
    });
  });

  describe('getSections', () => {
    it('calls the correct endpoint with resource ID', async () => {
      const mockSections = [
        {id: 1, resource_id: 5, kind: 'chapter', position: 1, title: 'Wrestling', description: 'Unarmed combat'},
      ];
      mockedGet.mockResolvedValue({data: mockSections});

      const result = await getSections(5);

      expect(mockedGet).toHaveBeenCalledWith(
        '/api/resources/5/sections',
      );
      expect(result).toEqual(mockSections);
    });
  });

  describe('getItems', () => {
    it('calls the correct endpoint with section ID', async () => {
      const mockItems = [
        {id: 1, section_id: 3, kind: 'technique', title: 'Zornhau', description: 'Wrath Strike', position: 1},
      ];
      mockedGet.mockResolvedValue({data: mockItems});

      const result = await getItems(3);

      expect(mockedGet).toHaveBeenCalledWith(
        '/api/sections/3/items',
      );
      expect(result).toEqual(mockItems);
    });
  });
});
