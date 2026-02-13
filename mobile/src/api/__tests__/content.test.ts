import apiClient from '../client';
import {getFightingBooks, getChapters, getTechniques} from '../content';

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

  describe('getFightingBooks', () => {
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

      const result = await getFightingBooks();

      expect(mockedGet).toHaveBeenCalledWith('/api/fighting-books', {
        params: {page: undefined, page_size: undefined},
      });
      expect(result).toEqual(mockResponse.data);
    });

    it('passes pagination params correctly', async () => {
      const mockResponse = {
        data: {
          data: [{id: 1, title: 'Test Book'}],
          page: 2,
          page_size: 10,
          total_count: 15,
          total_pages: 2,
        },
      };
      mockedGet.mockResolvedValue(mockResponse);

      const result = await getFightingBooks({page: 2, page_size: 10});

      expect(mockedGet).toHaveBeenCalledWith('/api/fighting-books', {
        params: {page: 2, page_size: 10},
      });
      expect(result).toEqual(mockResponse.data);
    });
  });

  describe('getChapters', () => {
    it('calls the correct endpoint with book ID', async () => {
      const mockChapters = [
        {id: 1, fighting_book_id: 5, chapter_number: 1, title: 'Wrestling'},
      ];
      mockedGet.mockResolvedValue({data: mockChapters});

      const result = await getChapters(5);

      expect(mockedGet).toHaveBeenCalledWith(
        '/api/fighting-books/5/chapters',
      );
      expect(result).toEqual(mockChapters);
    });
  });

  describe('getTechniques', () => {
    it('calls the correct endpoint with chapter ID', async () => {
      const mockTechniques = [
        {id: 1, chapter_id: 3, name: 'Zornhau', order_in_chapter: 1},
      ];
      mockedGet.mockResolvedValue({data: mockTechniques});

      const result = await getTechniques(3);

      expect(mockedGet).toHaveBeenCalledWith(
        '/api/chapters/3/techniques',
      );
      expect(result).toEqual(mockTechniques);
    });
  });
});
