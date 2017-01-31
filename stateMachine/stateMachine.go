package stateMachine

type StateMachine struct {
    States map[string]State
    stateChangeChannel chan string
    currentState State
    shutdownChannel chan bool
}

func NewStateMachine() *StateMachine {
    return &StateMachine {
        States: make(map[string]State),
        stateChangeChannel: make(chan string),
        shutdownChannel: make(chan bool),
    }
}

func (s *StateMachine) AddState(name string, state State) {
    s.States[name] = state
}

func (s *StateMachine) ChangeState(stateKey string) {
    s.stateChangeChannel <- stateKey
}

func (s *StateMachine) Start(startState string) {
    s.currentState = s.States[startState]

    if s.currentState != nil {
        s.currentState.OnEnter(nil)

        go s.updateLoop()
    }
}

func (s *StateMachine) ShutDown() {
    s.shutdownChannel <- true
}

func (s *StateMachine) updateLoop() {
    for {
        select {
            case stateName := <-s.stateChangeChannel:
                s.changeState(stateName)
            case <-s.shutdownChannel:
                break
            default:
                s.currentState.Update()
        }
    }
}

func (s *StateMachine) changeState(stateKey string) {
    newState := s.States[stateKey]

    if newState != nil {
        data := s.currentState.OnExit()
        s.currentState = newState
        s.currentState.OnEnter(data)
    }
}

type State interface {
    OnEnter(interface{})
    Update()
    OnExit() interface{}
}
