import React from 'react';
import {render, fireEvent, waitFor} from '@testing-library/react-native';
import {QueryClient, QueryClientProvider} from '@tanstack/react-query';
import TechniquesScreen from '../TechniquesScreen';

const mockNavigate = jest.fn();
const mockGoBack = jest.fn();
jest.mock('@react-navigation/native', () => ({
  useNavigation: () => ({navigate: mockNavigate, goBack: mockGoBack}),
  useRoute: () => ({
    params: {sectionId: 3, sectionTitle: 'Longsword'},
  }),
}));

const mockGetItems = jest.fn();
jest.mock('@/api/content', () => ({
  getItems: (...args: unknown[]) => mockGetItems(...args),
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

describe('TechniquesScreen', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders the section title in the header', async () => {
    mockGetItems.mockResolvedValue([]);

    const {getByText} = render(<TechniquesScreen />, {
      wrapper: createWrapper(),
    });

    expect(getByText('Longsword')).toBeTruthy();
    expect(getByText('Techniques')).toBeTruthy();
  });

  it('renders item list', async () => {
    mockGetItems.mockResolvedValue([
      {
        id: 1,
        section_id: 3,
        kind: 'technique',
        title: 'Posta di Donna',
        description: "The Woman's Guard",
        position: 1,
        attributes: {instructions: 'Hold the sword near your right shoulder'},
      },
      {
        id: 2,
        section_id: 3,
        kind: 'technique',
        title: 'Zornhau',
        description: 'The Wrath Strike',
        position: 2,
        attributes: {instructions: 'Strike diagonally'},
      },
    ]);

    const {getByText} = render(<TechniquesScreen />, {
      wrapper: createWrapper(),
    });

    await waitFor(() => {
      expect(getByText('Posta di Donna')).toBeTruthy();
    });
    expect(getByText('Zornhau')).toBeTruthy();
  });

  it('navigates to TechniqueDetail when an item is pressed', async () => {
    const item = {
      id: 1,
      section_id: 3,
      kind: 'technique',
      title: 'Posta di Donna',
      description: "The Woman's Guard",
      position: 1,
      attributes: {instructions: 'Hold the sword near your right shoulder'},
    };
    mockGetItems.mockResolvedValue([item]);

    const {getByTestId} = render(<TechniquesScreen />, {
      wrapper: createWrapper(),
    });

    await waitFor(() => {
      expect(getByTestId('technique-button-1')).toBeTruthy();
    });

    fireEvent.press(getByTestId('technique-button-1'));
    expect(mockNavigate).toHaveBeenCalledWith('TechniqueDetail', {
      item,
    });
  });

  it('navigates back when back button is pressed', async () => {
    mockGetItems.mockResolvedValue([]);

    const {getByTestId} = render(<TechniquesScreen />, {
      wrapper: createWrapper(),
    });

    fireEvent.press(getByTestId('back-button'));
    expect(mockGoBack).toHaveBeenCalled();
  });

  it('shows empty message when no items', async () => {
    mockGetItems.mockResolvedValue([]);

    const {getByText} = render(<TechniquesScreen />, {
      wrapper: createWrapper(),
    });

    await waitFor(() => {
      expect(getByText('No items available yet.')).toBeTruthy();
    });
  });
});
