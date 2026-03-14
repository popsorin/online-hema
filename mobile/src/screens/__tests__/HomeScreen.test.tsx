import React from 'react';
import {render, fireEvent, waitFor} from '@testing-library/react-native';
import {QueryClient, QueryClientProvider} from '@tanstack/react-query';
import HomeScreen from '../HomeScreen';

const mockNavigate = jest.fn();
jest.mock('@react-navigation/native', () => ({
  useNavigation: () => ({navigate: mockNavigate}),
}));

const mockGetResources = jest.fn();
jest.mock('@/api/content', () => ({
  getResources: (...args: unknown[]) => mockGetResources(...args),
}));

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {retry: false},
    },
  });
  return ({children}: {children: React.ReactNode}) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('HomeScreen', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders the header title', async () => {
    mockGetResources.mockResolvedValue({
      data: [],
      page: 1,
      page_size: 20,
      total_count: 0,
      total_pages: 0,
    });

    const {getByText} = render(<HomeScreen />, {wrapper: createWrapper()});

    expect(getByText('HEMA Lessons')).toBeTruthy();
  });

  it('renders resource cards', async () => {
    mockGetResources.mockResolvedValue({
      data: [
        {
          id: 1,
          author_id: 2,
          title: 'Fior di Battaglia',
          description: 'A combat manual',
          publication_year: 1409,
          cover_image_url: undefined,
          author_name: 'Fiore dei Liberi',
        },
      ],
      page: 1,
      page_size: 20,
      total_count: 1,
      total_pages: 1,
    });

    const {getByText} = render(<HomeScreen />, {wrapper: createWrapper()});

    await waitFor(() => {
      expect(getByText('Fior di Battaglia')).toBeTruthy();
    });
    expect(getByText('Fiore dei Liberi')).toBeTruthy();
    expect(getByText('1409')).toBeTruthy();
  });

  it('navigates to Chapters when a resource card is pressed', async () => {
    mockGetResources.mockResolvedValue({
      data: [
        {
          id: 1,
          author_id: 2,
          title: 'Fior di Battaglia',
          description: 'A combat manual',
          publication_year: 1409,
          cover_image_url: undefined,
          author_name: 'Fiore dei Liberi',
        },
      ],
      page: 1,
      page_size: 20,
      total_count: 1,
      total_pages: 1,
    });

    const {getByTestId} = render(<HomeScreen />, {wrapper: createWrapper()});

    await waitFor(() => {
      expect(getByTestId('book-card-1')).toBeTruthy();
    });

    fireEvent.press(getByTestId('book-card-1'));
    expect(mockNavigate).toHaveBeenCalledWith('Chapters', {
      resourceId: 1,
      resourceTitle: 'Fior di Battaglia',
    });
  });

  it('shows empty message when no resources', async () => {
    mockGetResources.mockResolvedValue({
      data: [],
      page: 1,
      page_size: 20,
      total_count: 0,
      total_pages: 0,
    });

    const {getByText} = render(<HomeScreen />, {wrapper: createWrapper()});

    await waitFor(() => {
      expect(getByText('No resources available yet.')).toBeTruthy();
    });
  });
});
