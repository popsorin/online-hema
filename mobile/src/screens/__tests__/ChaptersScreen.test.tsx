import React from 'react';
import {render, fireEvent, waitFor} from '@testing-library/react-native';
import {QueryClient, QueryClientProvider} from '@tanstack/react-query';
import ChaptersScreen from '../ChaptersScreen';

const mockNavigate = jest.fn();
const mockGoBack = jest.fn();
jest.mock('@react-navigation/native', () => ({
  useNavigation: () => ({navigate: mockNavigate, goBack: mockGoBack}),
  useRoute: () => ({
    params: {resourceId: 1, resourceTitle: 'Fior di Battaglia'},
  }),
}));

const mockGetSections = jest.fn();
jest.mock('@/api/content', () => ({
  getSections: (...args: unknown[]) => mockGetSections(...args),
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

describe('ChaptersScreen', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders the resource title in the header', async () => {
    mockGetSections.mockResolvedValue([]);

    const {getByText} = render(<ChaptersScreen />, {
      wrapper: createWrapper(),
    });

    expect(getByText('Fior di Battaglia')).toBeTruthy();
    expect(getByText('Chapters')).toBeTruthy();
  });

  it('renders section list', async () => {
    mockGetSections.mockResolvedValue([
      {
        id: 1,
        resource_id: 1,
        kind: 'chapter',
        position: 1,
        title: 'Wrestling',
        description: 'Techniques for unarmed combat and grappling',
      },
      {
        id: 2,
        resource_id: 1,
        kind: 'chapter',
        position: 2,
        title: 'Dagger Combat',
        description: 'Fighting with the dagger in various situations',
      },
    ]);

    const {getByText} = render(<ChaptersScreen />, {
      wrapper: createWrapper(),
    });

    await waitFor(() => {
      expect(getByText('Wrestling')).toBeTruthy();
    });
    expect(getByText('Dagger Combat')).toBeTruthy();
  });

  it('navigates to Techniques when a section is pressed', async () => {
    mockGetSections.mockResolvedValue([
      {
        id: 3,
        resource_id: 1,
        kind: 'chapter',
        position: 3,
        title: 'Longsword',
        description: 'The art of fighting with the longsword',
      },
    ]);

    const {getByTestId} = render(<ChaptersScreen />, {
      wrapper: createWrapper(),
    });

    await waitFor(() => {
      expect(getByTestId('chapter-button-3')).toBeTruthy();
    });

    fireEvent.press(getByTestId('chapter-button-3'));
    expect(mockNavigate).toHaveBeenCalledWith('Techniques', {
      sectionId: 3,
      sectionTitle: 'Longsword',
    });
  });

  it('navigates back when back button is pressed', async () => {
    mockGetSections.mockResolvedValue([]);

    const {getByTestId} = render(<ChaptersScreen />, {
      wrapper: createWrapper(),
    });

    fireEvent.press(getByTestId('back-button'));
    expect(mockGoBack).toHaveBeenCalled();
  });

  it('shows empty message when no sections', async () => {
    mockGetSections.mockResolvedValue([]);

    const {getByText} = render(<ChaptersScreen />, {
      wrapper: createWrapper(),
    });

    await waitFor(() => {
      expect(getByText('No sections available yet.')).toBeTruthy();
    });
  });
});
